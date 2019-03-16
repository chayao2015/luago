package vm

import (
	. "luago/api"
)

// “SELF指令（iABC模式）把对象和方法拷贝到相邻的两个目标寄存器中。对象在寄存器中，索引由操作数B指定。
// 方法名在常量表里，索引由操作数C指定。目标寄存器索引由操作数A指定”
// R(A+1) := R(B); R(A) := R(B)[RK(C)]
func self(inst Instruction, vm LuaVM) {
	a, b, c := inst.ABC()
	a++
	b++

	vm.Copy(b, a+1)
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}

// CLOSURE指令（iBx模式）把当前Lua函数的子函数原型实例化为闭包，放入由
// 操作数A指定的寄存器中。子函数原型来自于当前函数原型的子函数原型表，索
// 引由操作数Bx指定
// R(A) := closure(KPROTO[Bx])
func closure(inst Instruction, vm LuaVM) {
	a, bx := inst.ABx()
	a++

	vm.LoadProto(bx)
	vm.Replace(a)
}

// CALL指令（iABC模式）调用Lua函数。其中被调函数位于寄存器中，索引由操
// 作数A指定。需要传递给被调函数的参数值也在寄存器中，紧挨着被调函数，数量
// 由操作数B指定。函数调用结束后，原先存放函数和参数值的寄存器会被返回值占
// 据，具体有多少个返回值则由操作数C指定
// R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
//TODO:
func call(inst Instruction, vm LuaVM) {
	a, b, c := inst.ABC()
	a++

	nArgs := pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	popResults(a, c, vm)
}

// “VARARG指令（iABC模式）把传递给当前函数的变长参数加载到连续多个寄存器中。
// 其中第一个寄存器的索引由操作数A指定，寄存器数量由操作数B指定，操作数C没有用。
// VARARG指令可以用如下伪代码表示。
// R(A), R(A+1), ..., R(A+B-2) = vararg
func vararg(inst Instruction, vm LuaVM) {
	a, b, _ := inst.ABC()
	a++

	if b != 1 { // b==0 or b>1
		vm.LoadVararg(b - 1)
		popResults(a, b, vm)
	}
}

// return R(A)(R(A+1), ... ,R(A+B-1))
func tailCall(inst Instruction, vm LuaVM) {
	a, b, _ := inst.ABC()
	a++

	// TODO: optimize tail call!
	c := 0
	nArgs := pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	popResults(a, c, vm)
}

// return R(A), ... ,R(A+B-2)
// “RETURN指令（iABC模式）把存放在连续多个寄存器里的值返回给主调函数。
// 其中第一个寄存器的索引由操作数A指定，寄存器数量由操作数B指定，操作数C没用”
func fReturn(inst Instruction, vm LuaVM) {
	a, b, _ := inst.ABC()
	a++

	if b == 1 {
		// no return values
	} else if b > 1 {
		vm.CheckStack(b - 1)
		for i := a; i <= a+b-2; i++ {
			vm.PushValue(i)
		}
	} else {
		fixStack(a, vm)
	}
}

//TODO:??/
//TFORCALL指令（iABC模式
// R(A+3), ... ,R(A+2+C) := R(A)(R(A+1), R(A+2));
func tForCall(inst Instruction, vm LuaVM) {
	a, _, c := inst.ABC()
	a++

	pushFuncAndArgs(a, 3, vm)
	vm.Call(2, c)
	popResults(a+3, c+1, vm)
}

func pushFuncAndArgs(a, b int, vm LuaVM) (nArgs int) {
	if b >= 1 {
		vm.CheckStack(b)
		for i := a; i < a+b; i++ {
			vm.PushValue(i)
		}
		return b - 1
	} else {
		fixStack(a, vm)
		return vm.GetTop() - vm.RegisterCount() - 1
	}
}

//TODO:
func fixStack(a int, vm LuaVM) {
	x := int(vm.ToInteger(-1))
	vm.Pop(1)

	vm.CheckStack(x - a)
	for i := a; i < x; i++ {
		vm.PushValue(i)
	}
	vm.Rotate(vm.RegisterCount()+1, x-a)
}

func popResults(a, c int, vm LuaVM) {
	if c == 1 {
		// no results
	} else if c > 1 {
		for i := a + c - 2; i >= a; i-- {
			vm.Replace(i)
		}
	} else {
		// leave results on stack
		vm.CheckStack(1)
		vm.PushInteger(int64(a))
	}
}
