package vm

import (
	. "luago/api"
)

//
// Lua的变量可以分为三类：局部变量在函数内部定义（本质上是函数调用帧里的寄存器），
// Upvalue是直接或间接外围函数定义的局部变量，
// 全局变量则是全局环境表的字段（通过隐藏的Upvalue，也就是_ENV进行访问）。

// Upvalue是非局部变量，换句话说，就是某外围函数中定义的局部变量
// 全局变量实际上是某个特殊的表的字段，而这个特殊的表正是我们全局环境
// 然后编译器会把全局变量的读写翻译成_ENV字段的读写，也就是说，全局变量实际上也是语法糖
//
//
//

// GETUPVAL指令（iABC模式），把当前闭包的某个Upvalue值拷贝到目标寄存
// 器中。其中目标寄存器的索引由操作数A指定，Upvalue索引由操作数B指定，操作数C没用
// R(A) := UpValue[B]
func getUpval(inst Instruction, vm LuaVM) {
	a, b, _ := inst.ABC()
	a++
	b++

	vm.Copy(LuaUpvalueIndex(b), a)
}

// SETUPVAL指令（iABC模式），使用寄存器中的值给当前闭包的Upvalue赋值。
// 其中寄存器索引由操作数A指定，Upvalue索引由操作数B指定，操作数C没用
// UpValue[B] := R(A)
func setUpval(inst Instruction, vm LuaVM) {
	a, b, _ := inst.ABC()
	a++
	b++
	vm.Copy(a, LuaUpvalueIndex(b))
}

// 如果当前闭包的某个Upvalue是表，则GETTABUP指令（iABC模式）可以根据
// 键从该表里取值，然后把值放入目标寄存器中。其中目标寄存器索引由操作数A
// 指定，Upvalue索引由操作数B指定，键（可能在寄存器中也可能在常量表中）索引
// 由操作数C指定。GETTABUP指令相当于GETUPVAL和GETTABLE这两条指令的
// 组合，不过前者的效率明显要高一些
// R(A) := UpValue[B][RK(C)]
func getTabUp(inst Instruction, vm LuaVM) {
	a, b, c := inst.ABC()
	a++
	b++

	vm.GetRK(c)
	vm.GetTable(LuaUpvalueIndex(b))
	vm.Replace(a)
}

// 如果当前闭包的某个Upvalue是表，则SETTABUP指令（iABC模式）可以根据
// 键往该表里写入值。其中Upvalue索引由操作数A指定，键和值可能在寄存器中也
// 可能在常量表中，索引分别由操作数B和C指定。和GETTABUP指令类似，
// SETTABUP指令相当于GETUPVAL和SETTABLE这两条指令的组合，不过一条指
// 令的效率要高一些
// UpValue[A][RK(B)] := RK(C)
func setTabUp(inst Instruction, vm LuaVM) {
	a, b, c := inst.ABC()
	a++

	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(LuaUpvalueIndex(a))
}
