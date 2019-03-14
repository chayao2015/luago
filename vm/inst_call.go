package vm

// import (
// 	. "luago/api"
// )

// CLOSURE指令（iBx模式）把当前Lua函数的子函数原型实例化为闭包，放入由
// 操作数A指定的寄存器中。子函数原型来自于当前函数原型的子函数原型表，索
// 引由操作数Bx指定
// R(A) := closure(KPROTO[Bx])
// func closure(inst Instruction, vm LuaVM) {
// 	a, bx := inst.ABx()
// 	a++

// 	vm.LoadProto(bx)
// 	vm.Replace(a)
// }

// CALL指令（iABC模式）调用Lua函数。其中被调函数位于寄存器中，索引由操
// 作数A指定。需要传递给被调函数的参数值也在寄存器中，紧挨着被调函数，数量
// 由操作数B指定。函数调用结束后，原先存放函数和参数值的寄存器会被返回值占
// 据，具体有多少个返回值则由操作数C指定
// R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
// func call(inst Instruction, vm LuaVM) {
// 	a, b, c := inst.ABC()
// 	//TODO:
// }
