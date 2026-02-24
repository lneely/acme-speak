install:V:
	go build -o $HOME/bin/Speak .
	./download-ggml-model.sh tiny

clean:V:
	rm -f $HOME/bin/Speak
