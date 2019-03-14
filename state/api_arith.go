package state

import (
	. "luago/api"
	"luago/number"
	"math"
)

type operator struct {
	integerFunc func(int64, int64) int64
	floatFunc   func(float64, float64) float64
}

var (
	iadd  = func(a, b int64) int64 { return a + b }
	fadd  = func(a, b float64) float64 { return a + b }
	isub  = func(a, b int64) int64 { return a - b }
	fsub  = func(a, b float64) float64 { return a - b }
	imul  = func(a, b int64) int64 { return a * b }
	fmul  = func(a, b float64) float64 { return a * b }
	imod  = number.IMod
	fmod  = number.FMod
	pow   = math.Pow
	div   = func(a, b float64) float64 { return a / b }
	iidiv = number.IFloorDiv
	fidiv = number.FFloorDiv
	band  = func(a, b int64) int64 { return a & b }
	bor   = func(a, b int64) int64 { return a | b }
	bxor  = func(a, b int64) int64 { return a ^ b }
	shl   = number.ShiftLeft
	shr   = number.ShiftRight
	iunm  = func(a, _ int64) int64 { return -a }
	funm  = func(a, _ float64) float64 { return -a }
	bnot  = func(a, _ int64) int64 { return ^a }
)

var operators = []operator{
	operator{iadd, fadd},
	operator{isub, fsub},
	operator{imul, fmul},
	operator{imod, fmod},
	operator{nil, pow},
	operator{nil, div},
	operator{iidiv, fidiv},
	operator{band, nil},
	operator{bor, nil},
	operator{bxor, nil},
	operator{shl, nil},
	operator{shr, nil},
	operator{iunm, funm},
	operator{bnot, nil},
}

// [-(2|1), +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_arith

//void lua_arith (lua_State *L, int op);
// 对栈顶的两个值（或者一个，比如取反）做一次数学或位操作。 其中，栈顶的那个值是第二个操作数。
// 它会弹出压入的值，并把结果放在栈顶。 这个函数遵循 Lua 对应的操作符运算规则 （即有可能触发元方法）。

// op 的值必须是下列常量中的一个：

// LUA_OPADD: 加法 (+)
// LUA_OPSUB: 减法 (-)
// LUA_OPMUL: 乘法 (*)
// LUA_OPDIV: 浮点除法 (/)
// LUA_OPIDIV: 向下取整的除法 (//)
// LUA_OPMOD: 取模 (%)
// LUA_OPPOW: 乘方 (^)
// LUA_OPUNM: 取负 (一元 -)
// LUA_OPBNOT: 按位取反 (~)
// LUA_OPBAND: 按位与 (&)
// LUA_OPBOR: 按位或 (|)
// LUA_OPBXOR: 按位异或 (~)
// LUA_OPSHL: 左移 (<<)
// LUA_OPSHR: 右移 (>>)
func (L *luaState) Arith(op ArithOp) {
	var a, b luaValue
	b = L.stack.pop()
	if op != LUA_OPUNM && op != LUA_OPBNOT {
		a = L.stack.pop()
	} else {
		a = b
	}

	operator := operators[op]
	if res := arith(a, b, operator); res != nil {
		L.stack.push(res)
	} else {
		panic("arithmetic error!")
	}
}

func arith(a, b luaValue, op operator) luaValue {
	if op.floatFunc == nil { // bitwise
		if x, ok := convertToInteger(a); ok {
			if y, ok := convertToInteger(b); ok {
				return op.integerFunc(x, y)
			}
		}
	} else { // arith
		if op.integerFunc != nil { // add,sub,mul,mod,idiv,unm
			if x, ok := a.(int64); ok {
				if y, ok := b.(int64); ok {
					return op.integerFunc(x, y)
				}
			}
		}
		if x, ok := convertToFloat(a); ok {
			if y, ok := convertToFloat(b); ok {
				return op.floatFunc(x, y)
			}
		}
	}
	return nil
}
