package state

import (
	. "luago/api"
)

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushnil
//将空值压栈
func (L *luaState) PushNil() {
	L.stack.push(nil)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushboolean
//把 b 作为一个布尔量压栈
func (L *luaState) PushBoolean(b bool) {
	L.stack.push(b)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushinteger
//把值为 n 的整数压栈
func (L *luaState) PushInteger(n int64) {
	L.stack.push(n)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushnumber
//把一个值为 n 的浮点数压栈
func (L *luaState) PushNumber(n float64) {
	L.stack.push(n)
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_pushstring
//将指针 s 指向的零结尾的字符串压栈。Lua 对这个字符串做一个内部副本（或是复用一个副本），
// 因此 s 处的内存在函数返回后，可以释放掉或是立刻重用于其它用途。
//返回内部副本的指针。如果 s 为 NULL，将 nil 压栈并返回 NULL
func (L *luaState) PushString(s string) {
	L.stack.push(s)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushcfunction
func (L *luaState) PushGoFunction(f GoFunction) {
	L.stack.push(newGoClosure(f))
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushglobaltable
//将全局环境压栈
func (L *luaState) PushGlobalTable() {
	global := L.registry.get(LUA_RIDX_GLOBALS)
	L.stack.push(global)
}
