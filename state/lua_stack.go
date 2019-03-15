package state

import (
	. "luago/api"
)

type luaStack struct {
	slots []luaValue
	top   int
	/* call info*/
	L       *luaState
	closure *closure
	varargs []luaValue
	openuvs map[int]*upvalue
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

func (S *luaStack) isValid(idx int) bool {
	if idx < LUA_REGISTRYINDEX {
		//索引小于注册表索引，说明是Upvalue伪索引，把它转成真实索引（从0开始）然后看它是否在有效范围之内
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := S.closure
		return c != nil && uvIdx < len(c.upvals)
	}
	if idx == LUA_REGISTRYINDEX {
		return true
	}
	absIdx := S.absIndex(idx)
	return absIdx > 0 && absIdx <= S.top
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
	if idx < LUA_REGISTRYINDEX {
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := S.closure
		if c == nil || uvIdx >= len(c.upvals) {
			return nil
		}
		return *(c.upvals[uvIdx].val)
	}

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
	if idx < LUA_REGISTRYINDEX {
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := S.closure
		if c != nil || uvIdx < len(c.upvals) {
			*(c.upvals[uvIdx].val) = val
		}
		return
	}

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
