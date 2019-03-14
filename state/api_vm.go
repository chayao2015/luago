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
}
func (L *luaState) LoadVararg(n int) {

}

func (L *luaState) LoadProto(idx int) {
	proto := L.stack.closure.proto.Protos[idx]
	c := newLuaClosure(proto)
	L.stack.push(c)
}
