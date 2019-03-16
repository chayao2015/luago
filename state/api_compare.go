package state

import (
	. "luago/api"
)

// [-0, +0, e]
// http://www.lua.org/manual/5.3/manual.html#lua_compare
// int lua_compare (lua_State *L, int index1, int index2, int op);
// 比较两个 Lua 值。 当索引 index1 处的值通过 op 和索引 index2 处的值做比较后条件满足，函数返回 1 。
// 这个函数遵循 Lua 对应的操作规则（即有可能触发元方法）。 反之，函数返回 0。 当任何一个索引无效时，函数也会返回 0 。

// op 值必须是下列常量中的一个：

// LUA_OPEQ: 相等比较 (==)
// LUA_OPLT: 小于比较 (<)
// LUA_OPLE: 小于等于比较 (<=)
func (L *luaState) Compare(idx1, idx2 int, op CompareOp) bool {
	if !L.stack.isValid(idx1) || !L.stack.isValid(idx2) {
		return false
	}

	a := L.stack.get(idx1)
	b := L.stack.get(idx2)
	switch op {
	case LUA_OPEQ:
		return eq(a, b, L)
	case LUA_OPLT:
		return lt(a, b, L)
	case LUA_OPLE:
		return le(a, b, L)
	default:
		panic("invalid compare op!")
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rawequal
//如果索引 index1 与索引 index2 处的值 本身相等（即不调用元方法），返回 1 。 否则返回 0 。 当任何一个索引无效时，也返回 0
func (L *luaState) RawEqual(idx1, idx2 int) bool {
	if !L.stack.isValid(idx1) || !L.stack.isValid(idx2) {
		return false
	}

	a := L.stack.get(idx1)
	b := L.stack.get(idx2)
	return eq(a, b, nil)
}

func eq(a, b luaValue, L *luaState) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case string:
		y, ok := b.(string)
		return ok && x == y
	case int64:
		switch y := b.(type) {
		case int64:
			return x == y
		case float64:
			return float64(x) == y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x == y
		case int64:
			return x == float64(y)
		default:
			return false
		}
	case *luaTable:
		if y, ok := b.(*luaTable); ok && x != y && L != nil {
			if res, ok := callMetamethod(x, y, "__eq", L); ok {
				return convertToBoolean(res)
			}
		}
		return a == b
	default:
		return a == b
	}
}

func lt(a, b luaValue, L *luaState) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x < y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x < y
		case float64:
			return float64(x) < y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x < y
		case int64:
			return x < float64(y)
		}
	}

	if res, ok := callMetamethod(a, b, "__lt", L); ok {
		return convertToBoolean(res)
	}
	panic("comparison error!")
}

func le(a, b luaValue, L *luaState) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x <= y
		case float64:
			return float64(x) <= y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x <= y
		case int64:
			return x <= float64(y)
		}
	}
	if res, ok := callMetamethod(a, b, "__le", L); ok {
		return convertToBoolean(res)
	} else if res, ok := callMetamethod(b, a, "__lt", L); ok {
		return !convertToBoolean(res)
	}
	panic("comparison error!")
}
