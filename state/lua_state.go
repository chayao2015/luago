package state

import (
	. "luago/api"
)

type luaState struct {
	registry *luaTable //register table
	stack    *luaStack // lua stack
}

func New() *luaState {
	registry := newLuaTable(0, 0)
	registry.put(LUA_RIDX_GLOBALS, newLuaTable(0, 0))
	L := &luaState{registry: registry}
	L.pushLuaStack(newLuaStack(LUA_MINSTACK, L))
	return L
}

func (L *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = L.stack
	L.stack = stack
}

func (L *luaState) popLuaStack() {
	stack := L.stack
	L.stack = stack.prev
	stack.prev = nil
}
