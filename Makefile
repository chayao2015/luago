.PHONY: run build clean

all: luac_win run

build:
	go build luago.go

run:
	go run luago.go

luac_win:
	./lua/bin/luac53.exe lua/hello_world.lua

clean:
	rm -rf luac.out
	go clean
