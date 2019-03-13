package vm

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

func (ins Instruction) ABC() (a, b, c int) {
	a = int(ins >> 6 & 0xFF)
	b = int(ins >> 14 & 0x1FF)
	c = int(ins >> 23 & 0x1FF)
	return
}

func (ins Instruction) ABx() (a, bx int) {
	a = int(ins >> 6 & 0xFF)
	bx = int(ins >> 14 & 0x1FF)
	return
}

func (ins Instruction) AsBx() (a, sbx int) {
	a, bx := ins.ABx()
	return a, bx - MAXARG_sBx
}

func (ins Instruction) Ax() int {
	return int(ins >> 6)
}

func (ins Instruction) Opcode() int {
	return int(ins & 0x3F)
}

func (ins Instruction) OpName() string {
	return opcodes[ins.Opcode()].name
}

func (self Instruction) OpMode() byte {
	return opcodes[self.Opcode()].opMode
}

func (self Instruction) BMode() byte {
	return opcodes[self.Opcode()].argBMode
}

func (self Instruction) CMode() byte {
	return opcodes[self.Opcode()].argCMode
}
