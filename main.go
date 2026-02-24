package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"9fans.net/go/acme"
)

func readLines(r io.Reader, ch chan string) {
	buf := make([]byte, 1)
	line := ""
	for {
		n, err := r.Read(buf)
		if n > 0 {
			if buf[0] == '\n' {
				// Newline = finalized output, send it
				if line != "" {
					ch <- line
					line = ""
				}
			} else if buf[0] == '\r' {
				// Carriage return = overwrite, keep accumulating
				line = ""
			} else {
				line += string(buf[0])
			}
		}
		if err != nil {
			if line != "" {
				ch <- line
			}
			break
		}
	}
}

const (
	stateFile = "/tmp/acme-voice-state"
	modelPath = "/home/lkn/src/acme-speak/models/ggml-tiny.bin"
)

func main() {
	// Kill any leftover whisper-stream processes
	exec.Command("pkill", "-f", "whisper-stream").Run()

	winid := os.Getenv("winid")
	
	// Check if already running
	if pidBytes, err := os.ReadFile(stateFile); err == nil {
		// Stop existing stream - use syscall to kill process group
		pid := strings.TrimSpace(string(pidBytes))
		pidInt, _ := strconv.Atoi(pid)
		if pidInt > 0 {
			// Kill the process and its children
			exec.Command("pkill", "-9", "-P", pid).Run()
			syscall.Kill(pidInt, syscall.SIGKILL)
		}
		os.Remove(stateFile)
		fmt.Fprintln(os.Stderr, "Stopped recording")
		return
	}

	// Start new stream - optimized for tiny model (faster, lower latency)
	cmd := exec.Command("whisper-stream", "-m", modelPath, "-t", "4", "--step", "500", "--length", "5000")
	
	// Read from both stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// Save PID
	os.WriteFile(stateFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644)

	// Open acme window if winid is set
	var win *acme.Win
	if winid != "" {
		id, _ := strconv.Atoi(winid)
		win, err = acme.Open(id, nil)
		if err != nil {
			log.Printf("Failed to open acme window: %v", err)
			win = nil
		}
	}

	// Process output from both pipes
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	
	// Merge stdout and stderr
	combined := make(chan string, 100)
	done := make(chan bool, 2)
	
	go func() {
		readLines(stdout, combined)
		done <- true
	}()
	go func() {
		readLines(stderr, combined)
		done <- true
	}()
	
	// Close combined when both readers finish
	go func() {
		<-done
		<-done
		close(combined)
	}()
	
	lastLine := ""
	for line := range combined {
		part := ansiRegex.ReplaceAllString(line, "")
		part = strings.TrimSpace(part)
		
		// Filter out debug/empty lines and hallucinations
		if part == "" || strings.HasPrefix(part, "whisper_") ||
			strings.HasPrefix(part, "main:") || strings.HasPrefix(part, "init:") ||
			part == "[BLANK_AUDIO]" || part == "[Start speaking]" ||
			strings.HasPrefix(part, "(") || strings.HasPrefix(part, "[Music]") {
			continue
		}

		// Only write if different from last line
		if part != lastLine {
			lastLine = part
			// Write to acme or stdout
			if win != nil {
				win.Write("body", []byte(part+" "))
			} else {
				fmt.Println(part)
			}
		}
	}

	cmd.Wait()
	os.Remove(stateFile)
	if win != nil {
		win.CloseFiles()
	}
}
