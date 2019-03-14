package vm

import (
	. "luago/api"
)

/* arith */

func add(i Instruction, vm LuaVM)  { binaryArith(i, vm, LUA_OPADD) }  // +
func sub(i Instruction, vm LuaVM)  { binaryArith(i, vm, LUA_OPSUB) }  // -
func mul(i Instruction, vm LuaVM)  { binaryArith(i, vm, LUA_OPMUL) }  // *
func mod(i Instruction, vm LuaVM)  { binaryArith(i, vm, LUA_OPMOD) }  // %
func pow(i Instruction, vm LuaVM)  { binaryArith(i, vm, LUA_OPPOW) }  // ^
func div(i Instruction, vm LuaVM)  { binaryArith(i, vm, LUA_OPDIV) }  // /
func idiv(i Instruction, vm LuaVM) { binaryArith(i, vm, LUA_OPIDIV) } // //
func band(i Instruction, vm LuaVM) { binaryArith(i, vm, LUA_OPBAND) } // &
func bor(i Instruction, vm LuaVM)  { binaryArith(i, vm, LUA_OPBOR) }  // |
func bxor(i Instruction, vm LuaVM) { binaryArith(i, vm, LUA_OPBXOR) } // ~
func shl(i Instruction, vm LuaVM)  { binaryArith(i, vm, LUA_OPSHL) }  // <<
func shr(i Instruction, vm LuaVM)  { binaryArith(i, vm, LUA_OPSHR) }  // >>
func unm(i Instruction, vm LuaVM)  { unaryArith(i, vm, LUA_OPUNM) }   // -
func bnot(i Instruction, vm LuaVM) { unaryArith(i, vm, LUA_OPBNOT) }  // ~

// R(A) := RK(B) op RK(C)
// 二元算术运算指令（iABC模式），对两个寄存器或常量值（索引由操作数B和C
// 指定）进行运算，将结果放入另一个寄存器（索引由操作数A指定）
func binaryArith(inst Instruction, vm LuaVM, op ArithOp) {
	a, b, c := inst.ABC()
	a++

	vm.GetRK(b)
	vm.GetRK(c)
	vm.Arith(op)
	vm.Replace(a)
}

// R(A) := op R(B)
// 一元算术运算指令（iABC模式），对操作数B所指定的寄存器里的值进行运
// 算，然后把结果放入操作数A所指定的寄存器中，操作数C没用
func unaryArith(inst Instruction, vm LuaVM, op ArithOp) {
	a, b, _ := inst.ABC()
	a++
	b++

	vm.PushValue(b)
	vm.Arith(op)
	vm.Replace(a)
}

// LEN指令（iABC模式）进行的操作和一元算术运算指令类似
// R(A) := length of R(B)
func length(inst Instruction, vm LuaVM) {
	a, b, _ := inst.ABC()
	a++
	b++

	vm.Len(b)
	vm.Replace(a)
}

// R(A) := R(B).. ... ..R(C)
// CONCAT指令（iABC模式），将连续n个寄存器（起止索引分别由操作数B和C
// 指定）里的值拼接，将结果放入另一个寄存器（索引由操作数A指定)
func concat(inst Instruction, vm LuaVM) {
	a, b, c := inst.ABC()
	a++
	b++
	c++

	n := c - b + 1
	vm.CheckStack(n)
	for i := b; i <= c; i++ {
		vm.PushValue(i)
	}
	vm.Concat(n)
	vm.Replace(a)
}

/* compare */
func eq(inst Instruction, vm LuaVM) { compare(inst, vm, LUA_OPEQ) } // ==
func lt(inst Instruction, vm LuaVM) { compare(inst, vm, LUA_OPLT) } // <
func le(inst Instruction, vm LuaVM) { compare(inst, vm, LUA_OPLE) } // <=

// if ((RK(B) op RK(C)) ~= A) then pc++
// 比较指令（iABC模式），比较寄存器或常量表里的两个值（索引分别由操作数
// B和C指定），如果比较结果和操作数A（转换为布尔值）匹配，则跳过下一条指令。
// 比较指令不改变寄存器状态
func compare(inst Instruction, vm LuaVM, op CompareOp) {
	a, b, c := inst.ABC()

	vm.GetRK(b)
	vm.GetRK(c)
	if vm.Compare(-2, -1, op) != (a != 0) {
		vm.AddPC(1)
	}
	vm.Pop(2)
}

/* logical */

// R(A) := not R(B)
// NOT指令（iABC模式）进行的操作和一元算术运算指令类似
func not(inst Instruction, vm LuaVM) {
	a, b, _ := inst.ABC()
	a++
	b++

	vm.PushBoolean(!vm.ToBoolean(b))
	vm.Replace(a)
}

// if not (R(A) <=> C) then pc++
// TEST指令（iABC模式），判断寄存器A（索引由操作数A指定）中的值转换为布
// 尔值之后是否和操作数C表示的布尔值一致，如果一致，则跳过下一条指令。
// TEST指令不使用操作数B，也不改变寄存器状态
func test(inst Instruction, vm LuaVM) {
	a, _, c := inst.ABC()
	a++

	if vm.ToBoolean(a) != (c != 0) {
		vm.AddPC(1)
	}
}

// if (R(B) <=> C) then R(A) := R(B) else pc++
//TESTSET指令（iABC模式），判断寄存器B（索引由操作数B指定）中的值转换
// 为布尔值之后是否和操作数C表示的布尔值一致，如果一致则将寄存器B中的值复
// 制到寄存器A（索引由操作数A指定）中，否则跳过下一条指令
func testSet(inst Instruction, vm LuaVM) {
	a, b, c := inst.ABC()
	a++
	b++

	if vm.ToBoolean(b) == (c != 0) {
		vm.Copy(b, a)
	} else {
		vm.AddPC(1)
	}
}
