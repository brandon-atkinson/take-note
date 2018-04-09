tn: tn.rkt
	raco exe $^

clean: 
	rm tn

install: tn
	cp $^ ~/bin
