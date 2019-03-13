.PHONY: run build clean

all: run

build:
	go build luago.go

run:
	go run luago.go

luac:
	./lua/bin/luac53.exe lua/hello_world.lua

clean:
	rm -rf luac.out
	go clean
