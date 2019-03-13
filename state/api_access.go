package state

import (
	"fmt"
	. "luago/api"
)

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_typename
func (L *luaState) TypeName(tp LuaType) string {
	switch tp {
	case LUA_TNONE:
		return "no value"
	case LUA_TNIL:
		return "nil"
	case LUA_TBOOLEAN:
		return "boolean"
	case LUA_TNUMBER:
		return "number"
	case LUA_TSTRING:
		return "string"
	case LUA_TTABLE:
		return "table"
	case LUA_TFUNCTION:
		return "function"
	case LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_type
func (L *luaState) Type(idx int) LuaType {
	if L.stack.isValid(idx) {
		val := L.stack.get(idx)
		return typeOf(val)
	}
	return LUA_TNONE
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnone
func (L *luaState) IsNone(idx int) bool {
	return L.Type(idx) == LUA_TNONE
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnil
func (L *luaState) IsNil(idx int) bool {
	return L.Type(idx) == LUA_TNIL
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnoneornil
func (L *luaState) IsNoneOrNil(idx int) bool {
	return L.Type(idx) <= LUA_TNIL
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isboolean
func (L *luaState) IsBoolean(idx int) bool {
	return L.Type(idx) == LUA_TBOOLEAN
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_istable
func (L *luaState) IsTable(idx int) bool {
	return L.Type(idx) == LUA_TTABLE
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isfunction
func (L *luaState) IsFunction(idx int) bool {
	return L.Type(idx) == LUA_TFUNCTION
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isthread
func (L *luaState) IsThread(idx int) bool {
	return L.Type(idx) == LUA_TTHREAD
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isstring
func (L *luaState) IsString(idx int) bool {
	t := L.Type(idx)
	return t == LUA_TSTRING || t == LUA_TNUMBER
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnumber
func (L *luaState) IsNumber(idx int) bool {
	_, ok := L.ToNumberX(idx)
	return ok
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isinteger
func (L *luaState) IsInteger(idx int) bool {
	val := L.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_toboolean
func (L *luaState) ToBoolean(idx int) bool {
	val := L.stack.get(idx)
	return convertToBoolean(val)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tointeger
func (L *luaState) ToInteger(idx int) int64 {
	i, _ := L.ToIntegerX(idx)
	return i
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tointegerx
//将给定索引处的 Lua 值转换为带符号的整数类型 lua_Integer。
//这个 Lua 值必须是一个整数，或是一个可以被转换为整数 （参见 §3.4.3）的数字或字符串； 否则，lua_tointegerx 返回 0 。
//如果 isnum 不是 NULL， *isnum 会被设为操作是否成功
func (L *luaState) ToIntegerX(idx int) (int64, bool) {
	val := L.stack.get(idx)
	i, ok := val.(int64)
	return i, ok
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tonumber
func (L *luaState) ToNumber(idx int) float64 {
	n, _ := L.ToNumberX(idx)
	return n
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tonumberx
//把给定索引处的 Lua 值转换为 lua_Number 这样一个 C 类型 （参见 lua_Number ）。
// 这个 Lua 值必须是一个数字或是一个可转换为数字的字符串 （参见 §3.4.3）； 否则， lua_tonumberx 返回 0 。
//如果 isnum 不是 NULL， *isnum 会被设为操作是否成功
func (L *luaState) ToNumberX(idx int) (float64, bool) {
	val := L.stack.get(idx)
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}

// [-0, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_tostring
func (L *luaState) ToString(idx int) string {
	s, _ := L.ToStringX(idx)
	return s
}

func (L *luaState) ToStringX(idx int) (string, bool) {
	val := L.stack.get(idx)

	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x) // todo
		L.stack.set(idx, s)
		return s, true
	default:
		return "", false
	}
}
