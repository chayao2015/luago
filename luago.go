package main

import(
	"io/ioutil"
	"luago/state"
)


func main() {
	chunkName := "luac.out"
	data, err := ioutil.ReadFile(chunkName)
	if err != nil {
		panic(err)
	}

	L := state.New()
	L.Load(data, chunkName, "b")
	L.Call(0, 0)
}