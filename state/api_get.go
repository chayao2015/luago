package state

import (
	. "luago/api"
)

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_newtable
// void lua_newtable (lua_State *L);
// 创建一张空表，并将其压栈。 它等价于 lua_createtable(L, 0, 0)
func (L *luaState) NewTable() {
	L.CreateTable(0, 0)
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_createtable
// void lua_createtable (lua_State *L, int narr, int nrec);
// 创建一张新的空表压栈。 参数 narr 建议了这张表作为序列使用时会有多少个元素； 参数 nrec 建议了这张表可能拥有多少序列之外的元素。
// Lua 会使用这些建议来预分配这张新表。 如果你知道这张表用途的更多信息，预分配可以提高性能。 否则，你可以使用函数 lua_newtable
func (L *luaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	L.stack.push(t)
}

// [-1, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_gettable
// int lua_gettable (lua_State *L, int index);
// 把 t[k] 的值压栈， 这里的 t 是指索引指向的值， 而 k 则是栈顶放的值。
// 这个函数会弹出堆栈上的键，把结果放在栈上相同位置。 和在 Lua 中一样， 这个函数可能触发对应 "index" 事件的元方法 （参见 §2.4 ）。
// 返回压入值的类型
func (L *luaState) GetTable(idx int) LuaType {
	t := L.stack.get(idx)
	k := L.stack.pop()
	return L.getTable(t, k, false)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_getfield
// int lua_getfield (lua_State *L, int index, const char *k);
// 把 t[k] 的值压栈， 这里的 t 是索引指向的值。 在 Lua 中，这个函数可能触发对应 "index" 事件对应的元方法 （参见 §2.4 ）。
// 函数将返回压入值的类型
func (L *luaState) GetField(idx int, k string) LuaType {
	t := L.stack.get(idx)
	return L.getTable(t, k, false)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_geti
// int lua_geti (lua_State *L, int index, lua_Integer i);
// 把 t[i] 的值压栈， 这里的 t 指给定的索引指代的值。 和在 Lua 里一样，这个函数可能会触发 "index" 事件的元方法 （参见 §2.4）。
// 返回压入值的类型
func (L *luaState) GetI(idx int, i int64) LuaType {
	t := L.stack.get(idx)
	return L.getTable(t, i, false)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_getglobal
//把全局变量 name 里的值压栈，返回该值的类型
func (L *luaState) GetGlobal(name string) LuaType {
	t := L.registry.get(LUA_RIDX_GLOBALS)
	return L.getTable(t, name, false)
}

// [-1, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawget
//类似于 lua_gettable ， 但是作一次直接访问（不触发元方法）
func (L *luaState) RawGet(idx int) LuaType {
	t := L.stack.get(idx)
	k := L.stack.pop()
	return L.getTable(t, k, true)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawgeti
// 把 t[n] 的值压栈， 这里的 t 是指给定索引处的表。 这是一次直接访问；就是说，它不会触发元方法。返回入栈值的类型
func (L *luaState) RawGetI(idx int, i int64) LuaType {
	t := L.stack.get(idx)
	return L.getTable(t, i, true)
}

// [-0, +(0|1), –]
// http://www.lua.org/manual/5.3/manual.html#lua_getmetatable
// int lua_getmetatable (lua_State *L, int index);
// 如果该索引处的值有元表，则将其元表压栈，返回 1 。 否则不会将任何东西入栈，返回 0
func (L *luaState) GetMetatable(idx int) bool {
	val := L.stack.get(idx)

	if mt := getMetatable(val, L); mt != nil {
		L.stack.push(mt)
		return true
	}
	return false
}

//
// raw true，表示需要忽略元方法
func (L *luaState) getTable(t, k luaValue, raw bool) LuaType {
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		if raw || v != nil || !tbl.hasMetafield("__index") {
			L.stack.push(v)
			return typeOf(v)
		}
	}

	if !raw {
		if mf := getMetafield(t, "__index", L); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				return L.getTable(x, k, false)
			case *closure:
				L.stack.push(mf)
				L.stack.push(t)
				L.stack.push(k)
				L.Call(2, 1)
				v := L.stack.get(-1)
				return typeOf(v)
			}
		}
	}
	panic("not a table!")
}
