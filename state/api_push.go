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
	L.stack.push(newGoClosure(f, 0))
}

// [-n, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_pushcclosure
// void lua_pushcclosure (lua_State *L, lua_CFunction fn, int n);
// 把一个新的 C 闭包压栈。
// 当创建了一个 C 函数后， 你可以给它关联一些值， 这就是在创建一个 C 闭包（参见 §4.4）； 接下来无论函数何时被调用，这些值都可以被这个函数访问到。
//  为了将一些值关联到一个 C 函数上， 首先这些值需要先被压入堆栈（如果有多个值，第一个先压）。
//   接下来调用 lua_pushcclosure 来创建出闭包并把这个 C 函数压到栈上。 参数 n 告之函数有多少个值需要关联到函数上。
//   lua_pushcclosure 也会把这些值从栈上弹出。n 的最大值是 255 。
// 当 n 为零时， 这个函数将创建出一个 轻量 C 函数， 它就是一个指向 C 函数的指针。 这种情况下，不可能抛出内存错误
func (L *luaState) PushGoClosure(f GoFunction, n int) {
	c := newGoClosure(f, n)
	for i := n; i > 0; i-- {
		val := L.stack.pop()
		c.upvals[n-1] = &upvalue{&val}
	}
	L.stack.push(c)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushglobaltable
//将全局环境压栈
func (L *luaState) PushGlobalTable() {
	global := L.registry.get(LUA_RIDX_GLOBALS)
	L.stack.push(global)
}
