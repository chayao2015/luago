package vm

import (
	. "luago/api"
)

// MOVE指令（iABC模式）把源寄存器（索引由操作数B指定）里的值移动到目标
// 寄存器（索引由操作数A指定）里
// R(A) := R(B)
func move(inst Instruction, vm LuaVM) {
	a, b, _ := inst.ABC()
	a++
	b++

	vm.Copy(b, a)
}

// JMP指令（iAsBx模式）执行无条件跳转。该指令往往和后面要介绍的TEST等
// 指令配合使用，但是也可能会单独出现
// JMP指令的操作数A和Upvalue有关
// pc+=sBx; if (A) close all upvalues >= R(A - 1)
func jmp(inst Instruction, vm LuaVM) {
	a, sBx := inst.AsBx()

	vm.AddPC(sBx)
	if a != 0 {
		vm.CloseUpvalues(a)
	}
}
