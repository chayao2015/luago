package vm

import (
	. "luago/api"
)

// R(A), R(A+1), ..., R(A+B) := nil
//LOADNIL指令（iABC模式）用于给连续n个寄存器放置nil值。
// 寄存器的起始索引由操作数A指定，寄存器数量则由操作数B指定，操作数C没有用
func loadNil(inst Instruction, vm LuaVM) {
	a, b, _ := inst.ABC()
	a++

	vm.PushNil()
	for i := a; i <= a+b; i++ {
		vm.Copy(-1, i)
	}
	vm.Pop(1)
}

// R(A) := (bool)B; if (C) pc++
// LOADBOOL指令（iABC模式）给单个寄存器设置布尔值。
// 寄存器索引由操作数A指定，布尔值由寄存器B指定（0代表false，非0代表true），
// 如果寄存器C非0则跳过下一条指令
func loadBool(inst Instruction, vm LuaVM) {
	a, b, c := inst.ABC()
	a++

	vm.PushBoolean(b != 0)
	vm.Replace(a)
	if c != 0 {
		vm.AddPC(1)
	}
}

// R(A) := Kst(Bx)
// LOADK指令（iABx模式）将常量表里的某个常量加载到指定寄存器，寄存器
// 索引由操作数A指定，常量表索引由操作数Bx指定。
// 如果用Kst（N）表示常量表中的第N个常量
func loadK(inst Instruction, vm LuaVM) {
	a, bx := inst.ABx()
	a++

	vm.GetConst(bx)
	vm.Replace(a)
}

// R(A) := Kst(extra arg)
// LOADKX指令（也是iABx模式）需要和EXTRAARG指令（iAx模式）搭配使用，
// 用后者的Ax操作数来指定常量索引。
// Ax操作数占26个比特，可以表达的最大无符号整数是67108864，可以满足大部分情况
func loadKx(inst Instruction, vm LuaVM) {
	a, _ := inst.ABx()
	a++
	ax := Instruction(vm.Fetch()).Ax()

	vm.GetConst(ax)
	vm.Replace(a)
}
