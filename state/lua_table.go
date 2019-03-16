package state

import (
	"luago/number"
	"math"
)

type luaTable struct {
	metatable *luaTable //元表
	arr       []luaValue
	mp        map[luaValue]luaValue
	/* used by next() */
	keys    map[luaValue]luaValue // 由于Go语言的map不保证遍历顺序（甚至同样内容的map，两次遍历的顺序也可能不一样），所以我们只能在遍历开始之前把所有的键都固定下来，保存在keys字段里
	lastKey luaValue
	changed bool
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

// 如果传入的键是nil，表示遍历开始，需要把所有的键都收集到keys
// 里。keys的键值对记录了表的键和下一个键的关系，因此keys字段初始化好之后，
// 直接根据传入参数取值并返回即可
func (T *luaTable) nextKey(key luaValue) luaValue {
	if T.keys == nil || key == nil {
		T.initKeys()
		T.changed = false
	}
	return T.keys[key]
}

func (T *luaTable) initKeys() {
	T.keys = make(map[luaValue]luaValue)
	var key luaValue
	for i, v := range T.arr {
		if v != nil {
			T.keys[key] = int64(i + 1) // key 为 nil 遍历开始
			key = int64(i + 1)
		}
	}

	for k, v := range T.mp {
		if v != nil {
			T.keys[key] = k
			key = k
		}
	}
	T.lastKey = key
}
