package vm

import (
	. "luago/api"
)

// Lua语言的for循环语句有两种形式：数值（Numerical）形式和通用（Generic）形
// 式。数值for循环用于按一定步长遍历某个范围内的数值，通用for循环主要用于遍历表

// 数值for循环需要借助两条指令来实现：FORPREP和 FORLOOP

// Lua编译器为了实现for循环，使用了三个
// 特殊的局部变量，这三个特殊局部变量的名字里都包含圆括号（属于非法标识
// 符），这样就避免了和程序中出现的普通变量重名的可能。由名字可知，这三个局
// 部变量分别存放数值、限制和步长，并且在循环开始之前就已经预先初始化好了
// （对应三条LOADK指令）

// 这三个特殊的变量正好对应前面伪代码中的R（A）、
// R（A+1）和R（A+2）这三个寄存器，我们自己在for循环里定义的变量i则对应
// R（A+3）寄存器。由此可知，FORPREP指令执行的操作其实就是在循环开始之前
// 预先给数值减去步长，然后跳转到FORLOOP指令正式开始循环

// FORLOOP指令则是先给数值加上步长，然后判断数值是否还在范围之内。如
// 果已经超出范围，则循环结束；若未超过范围则把数值拷贝给用户定义的局部变
// 量，然后跳转到循环体内部开始执行具体的代码块

// FORPREP R(A)-=R(A+2); pc+=sBx
func forPrep(inst Instruction, vm LuaVM) {
	a, sBx := inst.AsBx()
	a++

	if vm.Type(a) == LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a))
		vm.Replace(a)
	}
	if vm.Type(a+1) == LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a + 1))
		vm.Replace(a + 1)
	}
	if vm.Type(a+2) == LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a + 2))
		vm.Replace(a + 2)
	}

	vm.PushValue(a)
	vm.PushValue(a + 2)
	vm.Arith(LUA_OPSUB)
	vm.Replace(a)
	vm.AddPC(sBx)
}

//FORLOOP
// R(A)+=R(A+2);
// if R(A) <?= R(A+1) then {
//   pc+=sBx; R(A+3)=R(A)
// }
func forLoop(inst Instruction, vm LuaVM) {
	a, sBx := inst.AsBx()
	a++

	//R(A)+=R(A+2);
	vm.PushValue(a + 2)
	vm.PushValue(a)
	vm.Arith(LUA_OPADD)
	vm.Replace(a)

	isPositiveStep := vm.ToNumber(a+2) > 0
	if isPositiveStep && vm.Compare(a, a+1, LUA_OPLE) ||
		!isPositiveStep && vm.Compare(a+1, a, LUA_OPLE) {
		// pc+=sBx; R(A+3)=R(A)
		vm.AddPC(sBx)
		vm.Copy(a, a+3)
	}
}

//TFORLOOP指令（iAsBx模式）
// if R(A+1) ~= nil then {
//   R(A)=R(A+1); pc += sBx
// }
func tForLoop(inst Instruction, vm LuaVM) {
	a, sBx := inst.AsBx()
	a++

	if !vm.IsNil(a + 1) {
		vm.Copy(a+1, a)
		vm.AddPC(sBx)
	}
}
