package vm

import "luago/api"

// Lua虚拟机这里采用了一种叫作偏移二进制码（Offset Binary，也叫作Excess-K）的编码模
// 式。具体来说，如果把sBx解释成无符号整数时它的值是x，那么解释成有符号整数
// 时它的值就是x-K。那么K是什么呢？K取sBx所能表示的最大无符号整数值的一
// 半，也就是上面代码中的MAXARG_sBx

const MAXARG_Bx = 1<<18 - 1       // 262143
const MAXARG_sBx = MAXARG_Bx >> 1 // 131071

/*
 31       22       13       5    0
  +-------+^------+-^-----+-^-----
  |b=9bits |c=9bits |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |    bx=18bits    |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |   sbx=18bits    |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |    ax=26bits            |op=6|
  +-------+^------+-^-----+-^-----
 31      23      15       7      0
*/
type Instruction uint32

func (I Instruction) Opcode() int {
	return int(I & 0x3F)
}

func (I Instruction) ABC() (a, b, c int) {
	a = int(I >> 6 & 0xFF)
	b = int(I >> 23 & 0x1FF)
	c = int(I >> 14 & 0x1FF)
	return
}

func (I Instruction) ABx() (a, bx int) {
	a = int(I >> 6 & 0xFF)
	bx = int(I >> 14)
	return
}

func (I Instruction) AsBx() (a, sbx int) {
	a, bx := I.ABx()
	return a, bx - MAXARG_sBx
}

func (I Instruction) Ax() int {
	return int(I >> 6)
}

func (I Instruction) OpName() string {
	return opcodes[I.Opcode()].name
}

func (I Instruction) OpMode() byte {
	return opcodes[I.Opcode()].opMode
}

func (I Instruction) BMode() byte {
	return opcodes[I.Opcode()].argBMode
}

func (I Instruction) CMode() byte {
	return opcodes[I.Opcode()].argCMode
}

func (I Instruction) Execute(vm api.LuaVM) {
	action := opcodes[I.Opcode()].action
	if action != nil {
		action(I, vm)
	} else {
		panic(I.OpName())
	}
}
