package vm

import (
	. "luago/api"
)

/* number of list items to accumulate before a SETLIST instruction */
const LFIELDS_PER_FLUSH = 50

// R(A) := {} (size = B,C)
// NEWTABLE指令（iABC模式）创建空表，并将其放入指定寄存器。寄存器索引
// 由操作数A指定，表的初始数组容量和哈希表容量分别由操作数B和C指定
func newTable(inst Instruction, vm LuaVM) {
	a, b, c := inst.ABC()
	a++

	// 	Fb2int（）函数起到什么作用呢？因为NEWTABLE
	// 指令是iABC模式，操作数B和C只有9个比特，如果当作无符号整数的话，最大也不
	// 能超过512。但是我们在前面也提到过，因为表构造器便捷实用，所以Lua也经常被
	// 用来描述数据（类似JSON），如果有很大的数据需要写成表构造器，但是表的初始
	// 容量又不够大，就容易导致表频繁扩容从而影响数据加载效率。
	// 为了解决这个问题，NEWTABLE指令的B和C操作数使用了一种叫作浮点字
	// 节（Floating Point Byte）的编码方式。这种编码方式和浮点数的编码方式类似，只
	// 是仅用一个字节。具体来说，如果把某个字节用二进制写成eeeeexxx，那么当
	// eeeee==0时该字节表示的整数就是xxx，否则该字节表示的整数是（1xxx）*2^（eeeee-1）
	vm.CreateTable(Fb2int(b), Fb2int(c))
	vm.Replace(a)
}

// R(A) := R(B)[RK(C)]
// GETTABLE指令（iABC模式）根据键从表里取值，并放入目标寄存器中。其中
// 表位于寄存器中，索引由操作数B指定；键可能位于寄存器中，也可能在常量表
// 里，索引由操作数C指定；目标寄存器索引则由操作数A指定
func getTable(inst Instruction, vm LuaVM) {
	a, b, c := inst.ABC()
	a++
	b++

	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}

// R(A)[RK(B)] := RK(C)
// SETTABLE指令（iABC模式）根据键往表里赋值。其中表位于寄存器中，索引
// 由操作数A指定；键和值可能位于寄存器中，也可能在常量表里，索引分别由操作
// 数B和C指定
func setTable(inst Instruction, vm LuaVM) {
	a, b, c := inst.ABC()
	a++

	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(a)
}

// R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
// SETTABLE是通用指令，每次只处理一个键值对，具体操作交给表去处理，并
// 不关心实际写入的是表的哈希部分还是数组部分。SETLIST指令（iABC模式）则是
// 专门给数组准备的，用于按索引批量设置数组元素。其中数组位于寄存器中，索
// 引由操作数A指定；需要写入数组的一系列值也在寄存器中，紧挨着数组，数量由
// 操作数B指定；数组起始索引则由操作数C指定。
// 那么数组的索引到底是怎么计算的呢？这里的情况和GETTABLE指令有点类
// 似。因为C操作数只有9个比特，所以直接用它表示数组索引显然不够用。这里的
// 解决办法是让C操作数保存批次数，然后用批次数乘上批大小（对应伪代码中的
// FPF）就可以算出数组起始索引。以默认的批大小50为例，C操作数能表示的最大
// 索引就是扩大到了25600（50*512）。
// 但是如果数组长度大于25600呢？是不是后面的元素就只能用SETTABLE指
// 令设置了？也不是。这种情况下SETLIST指令后面会跟一条EXTRAARG指令，用
// 其Ax操作数来保存批次数。综上所述，如果指令的C操作数大于0，那么表示的是
// 批次数加1，否则，真正的批次数存放在后续的EXTRAARG指令里
func setList(inst Instruction, vm LuaVM) {
	a, b, c := inst.ABC()
	a++

	if c > 0 {
		c = c - 1
	} else {
		c = Instruction(vm.Fetch()).Ax()
	}

	vm.CheckStack(1)
	idx := int64(c * LFIELDS_PER_FLUSH)
	for j := 1; j <= b; j++ {
		idx++
		vm.PushValue(a + j)
		vm.SetI(a, idx)
	}
}
