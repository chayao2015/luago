package state

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_gettop
// 返回栈顶元素的索引。 因为索引是从 1 开始编号的， 所以这个结果等于栈上的元素个数； 特别指出，0 表示栈为空
func (L *luaState) GetTop() int {
	return L.stack.top
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_absindex
//将一个可接受的索引 idx 转换为绝对索引 （即，一个不依赖栈顶在哪的值）
func (L *luaState) AbsIndex(idx int) int {
	return L.stack.absIndex(idx)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_checkstack
//确保堆栈上至少有 n 个额外空位。 如果不能把堆栈扩展到相应的尺寸，函数返回假。
//失败的原因包括将把栈扩展到比固定最大尺寸还大 （至少是几千个元素）或分配内存失败。
//这个函数永远不会缩小堆栈； 如果堆栈已经比需要的大了，那么就保持原样
func (L *luaState) CheckStack(n int) bool {
	L.stack.check(n)
	return true
}

// [-n, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pop
//从栈中弹出 n 个元素
func (L *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		L.stack.pop()
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_copy
//从索引 fromidx 处复制一个值到一个有效索引 toidx 处，覆盖那里的原有值。 不会影响其它位置的值
func (L *luaState) Copy(fromIdx, toIdx int) {
	val := L.stack.get(fromIdx)
	L.stack.set(toIdx, val)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushvalue
//把栈上给定索引处的元素作一个副本压栈
func (L *luaState) PushValue(idx int) {
	val := L.stack.get(idx)
	L.stack.push(val)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_replace
//把栈顶元素放置到给定位置而不移动其它元素 （因此覆盖了那个位置处的值），然后将栈顶元素弹出
func (L *luaState) Replace(idx int) {
	val := L.stack.pop()
	L.stack.set(idx, val)
}

// [-1, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_insert
//把栈顶元素移动到指定的有效索引处， 依次移动这个索引之上的元素。
// 不要用伪索引来调用这个函数， 因为伪索引没有真正指向栈上的位置
func (L *luaState) Insert(idx int) {
	L.Rotate(idx, 1)
}

// [-1, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_remove
// 从给定有效索引处移除一个元素， 把这个索引之上的所有元素移下来填补上这个空隙。
// 不能用伪索引来调用这个函数，因为伪索引并不指向真实的栈上的位置
func (L *luaState) Remove(idx int) {
	L.Rotate(idx, -1)
	L.Pop(1)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_rotate
//把从 idx 开始到栈顶的元素轮转 n 个位置。 对于 n 为正数时，轮转方向是向栈顶的； 当 n 为负数时，向栈底方向轮转 -n 个位置。 n 的绝对值不可以比参于轮转的切片长度大
func (L *luaState) Rotate(idx, n int) {
	t := L.stack.top - 1           /* end of stack segment being rotated */
	p := L.stack.absIndex(idx) - 1 /* start of segment */
	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	L.stack.reverse(p, m)   /* reverse the prefix with length 'n' */
	L.stack.reverse(m+1, t) /* reverse the suffix */
	L.stack.reverse(p, t)   /* reverse the entire segment */
}

// [-?, +?, –]
// http://www.lua.org/manual/5.3/manual.html#lua_settop
//参数允许传入任何索引以及 0 。 它将把堆栈的栈顶设为这个索引。
//如果新的栈顶比原来的大， 超出部分的新元素将被填为 nil 。 如果 index 为 0 ， 把栈上所有元素移除
func (L *luaState) SetTop(idx int) {
	newTop := L.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow!")
	}
	n := L.stack.top - newTop
	if n > 0 {
		for i := 0; i < n; i++ {
			L.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			L.stack.push(nil)
		}
	}
}
