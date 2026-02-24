install:V:
	cp Speak $HOME/bin/Speak
	chmod +x $HOME/bin/Speak
	./download-ggml-model.sh base
clean:V:
	rm -f $HOME/bin/Speak
