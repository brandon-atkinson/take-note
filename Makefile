tn: tn.go
	go build $^

clean:
	rm tn

install: tn
	cp $^ ~/bin
