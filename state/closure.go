package state

import (
	"luago/binchunk"
)

type luaClosure struct {
	proto *binchunk.Prototype
}

func newLuaClosure(proto *binchunk.Prototype) *luaClosure {
	return &luaClosure{proto: proto}
}
