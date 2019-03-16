package state

import (
	. "luago/api"
)

// [-2, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_settable
// void lua_settable (lua_State *L, int index)
//做一个等价于 t[k] = v 的操作， 这里 t 是给出的索引处的值， v 是栈顶的那个值， k 是栈顶之下的值。
// 这个函数会将键和值都弹出栈。 跟在 Lua 中一样，这个函数可能触发一个 "newindex" 事件的元方法
func (L *luaState) SetTable(idx int) {
	t := L.stack.get(idx)
	v := L.stack.pop()
	k := L.stack.pop()
	L.setTable(t, k, v, false)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_setfield
//void lua_setfield (lua_State *L, int index, const char *k)
// 做一个等价于 t[k] = v 的操作， 这里 t 是给出的索引处的值， 而 v 是栈顶的那个值。
// 这个函数将把这个值弹出栈。 跟在 Lua 中一样，这个函数可能触发一个 "newindex" 事件的元方法
func (L *luaState) SetField(idx int, k string) {
	t := L.stack.get(idx)
	v := L.stack.pop()
	L.setTable(t, k, v, false)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_seti
// void lua_seti (lua_State *L, int index, lua_Integer n);
// 做一个等价于 t[n] = v 的操作， 这里 t 是给出的索引处的值， 而 v 是栈顶的那个值。
// 这个函数将把这个值弹出栈。 跟在 Lua 中一样，这个函数可能触发一个 "newindex" 事件的元方法
func (L *luaState) SetI(idx int, i int64) {
	t := L.stack.get(idx)
	v := L.stack.pop()
	L.setTable(t, i, v, false)
}

// [-1, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_setglobal
//从堆栈上弹出一个值，并将其设为全局变量 name 的新值
func (L *luaState) SetGlobal(name string) {
	t := L.registry.get(LUA_RIDX_GLOBALS)
	v := L.stack.pop()
	L.setTable(t, name, v, false)
}

// [-2, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawset
//类似于 lua_settable ， 但是是做一次直接赋值（不触发元方法）
func (L *luaState) RawSet(idx int) {
	t := L.stack.get(idx)
	v := L.stack.pop()
	k := L.stack.pop()
	L.setTable(t, k, v, true)
}

// [-1, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_rawseti
// 等价于 t[i] = v ， 这里的 t 是指给定索引处的表， 而 v 是栈顶的值。
// 这个函数会将值弹出栈。 赋值是直接的；即不会触发元方法。
func (L *luaState) RawSetI(idx int, i int64) {
	t := L.stack.get(idx)
	v := L.stack.pop()
	L.setTable(t, i, v, true)
}

// [-0, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_register
// 把 C 函数 f 设到全局变量 name
func (L *luaState) Register(name string, f GoFunction) {
	L.PushGoFunction(f)
	L.SetGlobal(name)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_setmetatable
//把一张表弹出栈，并将其设为给定索引处的值的元表
func (L *luaState) SetMetatable(idx int) {
	val := L.stack.get(idx)
	mtVal := L.stack.pop()
	if mtVal == nil {
		setMetatable(val, nil, L)
	} else if mt, ok := mtVal.(*luaTable); ok {
		setMetatable(val, mt, L)
	} else {
		panic("table expected!") // todo
	}
}

// t[k]=v
func (L *luaState) setTable(t, k, v luaValue, raw bool) {
	if tbl, ok := t.(*luaTable); ok {
		if raw || tbl.get(k) != nil || !tbl.hasMetafield("__newindex") {
			tbl.put(k, v)
			return
		}
	}

	if !raw {
		if mf := getMetafield(t, "__newindex", L); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				L.setTable(x, k, v, false)
				return
			case *closure:
				L.stack.push(mf)
				L.stack.push(t)
				L.stack.push(k)
				L.stack.push(v)
				L.Call(3, 0)
				return
			}
		}
	}

	panic("not a table!")
}
