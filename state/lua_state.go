package state

type luaState struct {
	stack *luaStack
}

func New() *luaState {
	return &luaState{
		stack: newLuaStack(20),
	}
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
