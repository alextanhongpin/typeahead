VERSION := $(shell git rev-parse --short HEAD)
start:
	go run cmd/main.go -i

prof:
	# go run cmd/main.go -cpu=profiling/$(VERSION)_cpu.out -mem=profiling/$(VERSION)_mem.out -i -in dict.gob
	go run cmd/main.go -cpu=profiling/$(VERSION)_cpu.out -mem=profiling/$(VERSION)_mem.out -i -source ./words.txt 

dictionary:
	curl -o words.txt https://raw.githubusercontent.com/dwyl/english-words/master/words.txt
