package state

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_len
// void lua_len (lua_State *L, int index);
// 返回给定索引的值的长度。 它等价于 Lua 中的 '#' 操作
// 它有可能触发 "length" 事件对应的元方法  结果压栈
func (L *luaState) Len(idx int) {
	val := L.stack.get(idx)
	if s, ok := val.(string); ok {
		L.stack.push(int64(len(s)))
	} else {
		panic("length error!")
	}
}

// [-n, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_concat
// void lua_concat (lua_State *L, int n);
// 连接栈顶的 n 个值， 然后将这些值出栈，并把结果放在栈顶。
// 如果 n 为 1 ，结果就是那个值放在栈上（即，函数什么都不做）；
//如果 n 为 0 ，结果是一个空串。 连接依照 Lua 中通常语义完成
func (L *luaState) Concat(n int) {
	if n == 0 {
		L.stack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			if L.IsString(-1) && L.IsString(-2) {
				s2 := L.ToString(-1)
				s1 := L.ToString(-2)
				L.stack.pop()
				L.stack.pop()
				L.stack.push(s1 + s2)
				continue
			}

			panic("concatenation error!")
		}
	}
}
