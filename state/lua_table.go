package state

import (
	"luago/number"
	"math"
)

type luaTable struct {
	metatable *luaTable //元表
	arr       []luaValue
	mp        map[luaValue]luaValue
}

func newLuaTable(nArr, nRec int) *luaTable {
	t := &luaTable{}
	if nArr > 0 {
		t.arr = make([]luaValue, 0, nArr)
	}
	if nRec > 0 {
		t.mp = make(map[luaValue]luaValue, nRec)
	}
	return t
}

func (T *luaTable) hasMetafield(filedName string) bool {
	return T.metatable != nil && T.metatable.get(filedName) != nil
}

func (T *luaTable) len() int {
	return len(T.arr)
}

func (T *luaTable) get(key luaValue) luaValue {
	key = floatToInteger(key)
	if idx, ok := key.(int64); ok {
		if idx >= 1 && idx <= int64(len(T.arr)) {
			return T.arr[idx-1]
		}
	}
	return T.mp[key]
}

func floatToInteger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
	}
	return key
}

func (T *luaTable) put(key, val luaValue) {
	if key == nil {
		panic("table index is nil!")
	}

	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is NaN!")
	}

	key = floatToInteger(key)

	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(T.arr))
		if idx <= arrLen {
			T.arr[idx-1] = val
			if idx == arrLen && val == nil {
				T.shrinkArray()
			}
			return
		}
		if idx == arrLen+1 {
			delete(T.mp, key)
			if val != nil {
				T.arr = append(T.arr, val)
				T.expandArray()
			}
			return
		}
	}

	if val != nil {
		if T.mp == nil {
			T.mp = make(map[luaValue]luaValue, 8)
		}
		T.mp[key] = val
	} else {
		delete(T.mp, key)
	}
}

func (T *luaTable) shrinkArray() {
	for i := len(T.arr) - 1; i >= 0; i-- {
		if T.arr[i] == nil {
			T.arr = T.arr[0:i]
		}
	}
}

func (T *luaTable) expandArray() {
	for idx := int64(len(T.arr)) + 1; true; idx++ {
		if val, found := T.mp[idx]; found {
			delete(T.mp, idx)
			T.arr = append(T.arr, val)
		} else {
			break
		}
	}
}
