package state

// [-2, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_settable
// void lua_settable (lua_State *L, int index)
//做一个等价于 t[k] = v 的操作， 这里 t 是给出的索引处的值， v 是栈顶的那个值， k 是栈顶之下的值。
// 这个函数会将键和值都弹出栈。 跟在 Lua 中一样，这个函数可能触发一个 "newindex" 事件的元方法
func (L *luaState) SetTable(idx int) {
	t := L.stack.get(idx)
	v := L.stack.pop()
	k := L.stack.pop()
	L.setTable(t, k, v)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_setfield
//void lua_setfield (lua_State *L, int index, const char *k)
// 做一个等价于 t[k] = v 的操作， 这里 t 是给出的索引处的值， 而 v 是栈顶的那个值。
// 这个函数将把这个值弹出栈。 跟在 Lua 中一样，这个函数可能触发一个 "newindex" 事件的元方法
func (L *luaState) SetField(idx int, k string) {
	t := L.stack.get(idx)
	v := L.stack.pop()
	L.setTable(t, k, v)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_seti
// void lua_seti (lua_State *L, int index, lua_Integer n);
// 做一个等价于 t[n] = v 的操作， 这里 t 是给出的索引处的值， 而 v 是栈顶的那个值。
// 这个函数将把这个值弹出栈。 跟在 Lua 中一样，这个函数可能触发一个 "newindex" 事件的元方法
func (L *luaState) SetI(idx int, i int64) {
	t := L.stack.get(idx)
	v := L.stack.pop()
	L.setTable(t, i, v)
}

func (L *luaState) setTable(t, k, v luaValue) {
	if tbl, ok := t.(*luaTable); ok {
		tbl.put(k, v)
		return
	}
	panic("not a table!")
}
