package state

import (
	"fmt"
	. "luago/api"
	"luago/number"
)

type luaValue interface{}

func typeOf(val luaValue) LuaType {
	switch val.(type) {
	case nil:
		return LUA_TNIL
	case bool:
		return LUA_TBOOLEAN
	case int64, float64:
		return LUA_TNUMBER
	case string:
		return LUA_TSTRING
	case *luaTable:
		return LUA_TTABLE
	case *closure:
		return LUA_TFUNCTION
	default:
		panic("TODO luaValue")
	}
}

func convertToBoolean(vale luaValue) bool {
	switch x := vale.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}

// http://www.lua.org/manual/5.3/manual.html#3.4.3
func convertToFloat(val luaValue) (float64, bool) {
	switch x := val.(type) {
	case int64:
		return float64(x), true
	case float64:
		return x, true
	case string:
		return number.ParseFloat(x)
	default:
		return 0, false
	}
}

// http://www.lua.org/manual/5.3/manual.html#3.4.3
func convertToInteger(val luaValue) (int64, bool) {
	switch x := val.(type) {
	case int64:
		return x, true
	case float64:
		return number.FloatToInteger(x)
	case string:
		return stringToInteger(x)
	default:
		return 0, false
	}
}

func stringToInteger(s string) (int64, bool) {
	if i, ok := number.ParseInteger(s); ok {
		return i, true
	}
	if f, ok := number.ParseFloat(s); ok {
		return number.FloatToInteger(f)
	}
	return 0, false
}

/* metatable */

// 先判断值是否是表，如果是，直接返回其元表字段即可；否则的话，
// 根据值的类型从注册表里取出与该类型关联的元表并返回，如果值没有元表与之关联，返回值就是nil
func getMetatable(val luaValue, L *luaState) *luaTable {
	if t, ok := val.(*luaTable); ok {
		return t.metatable
	}
	key := fmt.Sprintf("_MT%d", typeOf(val))
	if mt := L.registry.get(key); mt != nil {
		return mt.(*luaTable)
	}
	return nil
}

// 先判断值是否是表，如果是，直接修改其元表字段即可。否则的话，根据
// 变量类型把元表存储在注册表里，这样就达到了按类型共享元表的目的
// 虽然注册表也是一个普通的表，不过按照约定，下划线开头后跟大写字母的字段名是保
// 留给Lua实现使用的，所以我们使用了“_MT1”这样的字段名，以免和用户（通过
// API）放在注册表里的数据产生冲突。另外，如果传递给函数的元表是nil值，效果就相当于删除元表
func setMetatable(val luaValue, mt *luaTable, L *luaState) {
	if t, ok := val.(*luaTable); ok {
		t.metatable = mt
		return
	}
	key := fmt.Sprintf("_MT%d", typeOf(val))
	L.registry.put(key, val)
}

func getMetafield(val luaValue, fieldName string, L *luaState) luaValue {
	if mt := getMetatable(val, L); mt != nil {
		return mt.get(fieldName)
	}
	return nil
}

func callMetamethod(a, b luaValue, mmName string, L *luaState) (luaValue, bool) {
	var mm luaValue
	if mm = getMetafield(a, mmName, L); mm == nil {
		if mm = getMetafield(b, mmName, L); mm == nil {
			return nil, false
		}
	}
	L.stack.check(4)
	L.stack.push(mm)
	L.stack.push(a)
	L.stack.push(b)
	L.Call(2, 1)
	return L.stack.pop(), true
}
