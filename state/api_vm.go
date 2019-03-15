package state

func (L *luaState) PC() int {
	return L.stack.pc
}

func (L *luaState) AddPC(n int) {
	L.stack.pc += n
}

func (L *luaState) Fetch() uint32 {
	p := L.stack.closure.proto.Code[L.stack.pc]
	L.stack.pc++
	return p
}

func (L *luaState) GetConst(idx int) {
	c := L.stack.closure.proto.Constants[idx]
	L.stack.push(c)
}

// 传递给GetRK（）方法的参数实际上是iABC模式指令里的
// OpArgK类型参数。这种类型的参数一共占9个比特。如果最高位是
// 1，那么参数里存放的是常量表索引，把最高位去掉就可以得到索引值；否则最高
// 位是0，参数里存放的就是寄存器索引值
// Lua虚拟机指令操作
// 数里携带的寄存器索引是从0开始的，而Lua API里的栈索引是从1开始的，所以当
// 需要把寄存器索引当成栈索引使用时，要对寄存器索引加1
func (L *luaState) GetRK(rk int) {
	if rk > 0xFF { //constant
		L.GetConst(rk & 0xFF)
	} else {
		L.PushValue(rk + 1)
	}
}

func (L *luaState) RegisterCount() int {
	return int(L.stack.closure.proto.MaxStackSize)
}

func (L *luaState) LoadVararg(n int) {
	if n < 0 {
		n = len(L.stack.varargs)
	}

	L.stack.check(n)
	L.stack.pushN(L.stack.varargs, n)
}

// 根据函数原型里的Upvalue表来初始化闭包的Upvalue值。对于每个
// Upvalue，又有两种情况需要考虑：如果某一个Upvalue捕获的是当前函数的局部变
// 量（Instack==1），那么我们只要访问当前函数的局部变量即可；如果某一个
// Upvalue捕获的是更外围的函数中的局部变量（Instack==0），该Upvalue已经被当前
// 函数捕获，我们只要把该Upvalue传递给闭包即可。
// 对于第一种情况，如果Upvalue捕获的外围函数局部变量还在栈上，直接引用
// 即可，我们称这种Upvalue处于开放（Open）状态；反之，必须把变量的实际值保存
// 在其他地方，我们称这种Upvalue处于闭合（Closed）状态。为了能够在合适的时机
// （比如局部变量退出作用域时，详见10.3.5节）把处于开放状态的Upvalue闭合，需
// 要记录所有暂时还处于开放状态的Upvalue，我们把这些Upvalue记录在被捕获局
// 部变量所在的栈帧里。请读者打开luaStack.go文件（和closure.go文件在同一目录
// 下），给luaStack结构体添加openuvs字段。该字段是map类型，其中键是int类型，存
// 放局部变量的寄存器索引，值是Upvalue指针
func (L *luaState) LoadProto(idx int) {
	stk := L.stack
	subProto := stk.closure.proto.Protos[idx]
	c := newLuaClosure(subProto)
	stk.push(c)

	for i, uvInfo := range subProto.Upvalues {
		uvIdx := int(uvInfo.Idx)
		if uvInfo.Instack == 1 {
			if stk.openuvs == nil {
				stk.openuvs = map[int]*upvalue{}
			}

			if openuv, found := stk.openuvs[uvIdx]; found {
				c.upvals[i] = openuv
			} else {
				//********************//
				c.upvals[i] = &upvalue{&stk.slots[uvIdx]}
				stk.openuvs[uvIdx] = &upvalue{&stk.slots[i]}
			}
		} else {
			c.upvals[i] = stk.closure.upvals[uvIdx]
		}
	}
}

func (L *luaState) CloseUpvalues(a int) {
	for i, openuv := range L.stack.openuvs {
		if i >= a-1 {
			val := *openuv.val
			openuv.val = &val
			delete(L.stack.openuvs, i)
		}
	}
}
