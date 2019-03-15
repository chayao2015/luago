package api

type LuaVM interface {
	LuaState
	PC() int          // 返回当前PC（仅测试用） PC程序计数器 Program Counter)
	AddPC(n int)      // 修改PC（用于实现跳转指令）
	Fetch() uint32    //取出当前指令；将PC指向下一条指令
	GetConst(idx int) //从常量表将指定常量推入栈顶 LOADK和LOADKX这两个指令需要使用这个方法
	GetRK(rk int)     //将指定常量或栈值推入栈顶

	RegisterCount() int
	LoadVararg(n int)
	LoadProto(idx int)
	CloseUpvalues(a int)
}
