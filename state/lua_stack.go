package state

import (
	. "luago/api"
)

type luaStack struct {
	slots []luaValue
	top   int
	/* call info*/
	L       *luaState
	closure *Closure
	varargs []luaValue
	pc      int
	/* linked list*/
	prev *luaStack
}

func newLuaStack(size int, L *luaState) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
		L:     L,
	}
}

func (S *luaStack) absIndex(idx int) int {
	if idx >= 0 || idx <= LUA_REGISTRYINDEX {
		return idx
	}
	return idx + S.top + 1
}

func (S *luaStack) check(n int) {
	free := len(S.slots) - S.top
	for i := free; i < n; i++ {
		S.slots = append(S.slots, nil)
	}
}

func (S *luaStack) pop() luaValue {
	if S.top < 1 {
		panic("stack underflow!")
	}
	S.top--
	val := S.slots[S.top]
	S.slots[S.top] = nil
	return val
}

func (S *luaStack) push(val luaValue) {
	if S.top == len(S.slots) {
		panic("stack overflow!")
	}
	S.slots[S.top] = val
	S.top++
}

func (S *luaStack) popN(n int) []luaValue {
	vals := make([]luaValue, n)
	for i := n - 1; i >= 0; i-- {
		vals[i] = S.pop()
	}
	return vals
}

func (S *luaStack) pushN(vals []luaValue, n int) {
	nVals := len(vals)
	if n < 0 {
		n = nVals
	}

	for i := 0; i < n; i++ {
		if i < nVals {
			S.push(vals[i])
		} else {
			S.push(nil)
		}
	}
}

func (S *luaStack) get(idx int) luaValue {
	if idx == LUA_REGISTRYINDEX {
		return S.L.registry
	}

	absIdx := S.absIndex(idx)
	if absIdx > 0 && absIdx <= S.top {
		return S.slots[absIdx-1]
	}
	return nil
}

func (S *luaStack) set(idx int, val luaValue) {
	if idx == LUA_REGISTRYINDEX {
		S.L.registry = val.(*luaTable)
		return
	}

	absIdx := S.absIndex(idx)
	if absIdx > 0 && absIdx <= S.top {
		S.slots[absIdx-1] = val
	} else {
		panic("invalid index!")
	}
}

func (S *luaStack) reverse(from, to int) {
	slots := S.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}

func (S *luaStack) isValid(idx int) bool {
	if idx == LUA_REGISTRYINDEX {
		return true
	}
	absIdx := S.absIndex(idx)
	return absIdx > 0 && absIdx <= S.top
}
