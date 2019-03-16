package state

import (
	. "luago/api"
	"luago/binchunk"
	"luago/vm"
)

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_load
// int lua_load (lua_State *L,
// 	lua_Reader reader,
// 	void *data,
// 	const char *chunkname,
// 	const char *mode);
// 加载一段 Lua 代码块，但不运行它。 如果没有错误， lua_load 把一个编译好的代码块作为一个 Lua 函数压到栈顶。 否则，压入错误消息。

// lua_load 的返回值可以是：

// LUA_OK: 没有错误；
// LUA_ERRSYNTAX: 在预编译时碰到语法错误；
// LUA_ERRMEM: 内存分配错误；
// LUA_ERRGCMM: 在运行 __gc 元方法时出错了。 （这个错误和代码块加载过程无关，它是由垃圾收集器引发的。）
// lua_load 函数使用一个用户提供的 reader 函数来读取代码块（参见 lua_Reader ）。 data 参数会被传入 reader 函数。

// chunkname 这个参数可以赋予代码块一个名字， 这个名字被用于出错信息和调试信息（参见 §4.9）。

// lua_load 会自动检测代码块是文本的还是二进制的， 然后做对应的加载操作（参见程序 luac ）。 字符串 mode 的作用和函数 load 一致。 它还可以是 NULL 等价于字符串 "bt"。

// lua_load 的内部会使用栈， 因此 reader 函数必须永远在每次返回时保留栈的原样。

// 如果返回的函数有上值， 第一个上值会被设置为 保存在注册表（参见 §4.5） LUA_RIDX_GLOBALS 索引处的全局环境。 在加载主代码块时，这个上值是 _ENV 变量（参见 §2.2）。 其它上值均被初始化为 nil
func (L *luaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binchunk.Undump(chunk)
	c := newLuaClosure(proto)
	L.stack.push(c)
	// 如果需要，那么第一个Upvalue（对于主函数来说就是_ENV）会被初始化
	// 成全局环境，其他Upvalue会被初始化成nil
	if len(proto.Upvalues) > 0 {
		// 设置 _ENV
		env := L.registry.get(LUA_RIDX_GLOBALS)
		c.upvals[0] = &upvalue{&env}
	}
	return 0
}

// [-(nargs+1), +nresults, e]
// http://www.lua.org/manual/5.3/manual.html#lua_call
// void lua_call (lua_State *L, int nargs, int nresults);
// 调用一个函数。

// 要调用一个函数请遵循以下协议： 首先，要调用的函数应该被压入栈； 接着，把需要传递给这个函数的参数按正序压栈； 这是指第一个参数首先压栈。 最后调用一下 lua_call； nargs 是你压入栈的参数个数。
// 当函数调用完毕后，所有的参数以及函数本身都会出栈。 而函数的返回值这时则被压栈。
//  返回值的个数将被调整为 nresults 个， 除非 nresults 被设置成 LUA_MULTRET。 在这种情况下，所有的返回值都被压入堆栈中。
//  Lua 会保证返回值都放入栈空间中。 函数返回值将按正序压栈（第一个返回值首先压栈）， 因此在调用结束后，最后一个返回值将被放在栈顶。

// 被调用函数内发生的错误将（通过 longjmp ）一直上抛。

// 下面的例子中，这行 Lua 代码等价于在宿主程序中用 C 代码做一些工作：

//      a = f("how", t.x, 14)
// 这里是 C 里的代码：

//      lua_getglobal(L, "f");                  /* function to be called */
//      lua_pushliteral(L, "how");                       /* 1st argument */
//      lua_getglobal(L, "t");                    /* table to be indexed */
//      lua_getfield(L, -1, "x");        /* push result of t.x (2nd arg) */
//      lua_remove(L, -2);                  /* remove 't' from the stack */
//      lua_pushinteger(L, 14);                          /* 3rd argument */
//      lua_call(L, 3, 1);     /* call 'f' with 3 arguments and 1 result */
//      lua_setglobal(L, "a");                         /* set global 'a' */
// 注意上面这段代码是 平衡 的： 到了最后，堆栈恢复成原有的配置。 这是一种良好的编程习惯
func (L *luaState) Call(nArgs, nResults int) {
	val := L.stack.get(-(nArgs + 1))
	c, ok := val.(*closure)
	if !ok {
		//查找是否有 元方法 TODO:??
		if mf := getMetafield(val, "__call", L); mf != nil {
			if c, ok = mf.(*closure); ok {
				L.stack.push(val)
				L.Insert(-(nArgs + 2))
				nArgs++
			}
		}
	}

	if ok {
		if c.proto != nil {
			L.callLuaClosure(nArgs, nResults, c)
		} else {
			L.callGoClosure(nArgs, nResults, c)
		}
	} else {
		panic("not Closure!")
	}
}

// 先创建新的调用帧，然后把参数值从主调帧里弹出，推入被调帧。Go闭包直接从主调帧里弹出
// 扔掉即可。参数传递完毕之后，把被调帧推入调用栈，让它成为当前帧，然后直接
// 执行Go函数。执行完毕之后把被调帧从调用栈里弹出，这样主调帧就又成了当前
// 帧。最后（如果有必要），还需要把返回值从被调帧里弹出，推入主调帧（多退少补）
func (L *luaState) callGoClosure(nArgs, nResults int, c *closure) {
	// create new lua stack
	newStack := newLuaStack(nArgs+LUA_MINSTACK, L)
	newStack.closure = c

	// pass args, po func
	if nArgs > 0 {
		args := L.stack.popN(nArgs)
		newStack.pushN(args, nArgs)
	}
	L.stack.pop()

	// run closure
	L.pushLuaStack(newStack)
	n := c.goFunc(L)
	L.popLuaStack()

	// return results
	if nResults != 0 {
		results := newStack.popN(n)
		L.stack.pushN(results, nResults)
	}
}

//TODO:
func (L *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	// create new lua stack
	newStack := newLuaStack(nRegs+LUA_MINSTACK, L)
	newStack.closure = c

	// pass args, po func
	fundAndArgs := L.stack.popN(nArgs + 1)
	newStack.pushN(fundAndArgs[1:], nParams)
	newStack.top = nRegs
	if nArgs > nParams && isVararg {
		newStack.varargs = fundAndArgs[nParams+1:]
	}

	// run closure
	L.pushLuaStack(newStack)
	L.runLuaClosure()
	L.popLuaStack()

	// return results
	//把全部返回值从被调帧栈顶弹出， 然后根据期望的返回值数量多退少补，推入当前帧栈顶
	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		L.stack.check(len(results))
		L.stack.pushN(results, nResults)
	}
}

func (L *luaState) runLuaClosure() {
	for {
		inst := vm.Instruction(L.Fetch())
		inst.Execute(L)
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}
