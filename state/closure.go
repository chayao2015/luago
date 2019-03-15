package state

import (
	. "luago/api"
	"luago/binchunk"
)

type Closure struct {
	proto  *binchunk.Prototype
	goFunc GoFunction
}

func newLuaClosure(proto *binchunk.Prototype) *Closure {
	return &Closure{proto: proto}
}

func newGoClosure(f GoFunction) *Closure {
	return &Closure{goFunc: f}
}
