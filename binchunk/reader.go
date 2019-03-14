package binchunk

import "encoding/binary"
import "math"

type reader struct {
	data []byte
}

func (R *reader) readByte() byte {
	b := R.data[0]
	R.data = R.data[1:]
	return b
}

func (R *reader) readBytes(n uint) []byte {
	bytes := R.data[:n]
	R.data = R.data[n:]
	return bytes
}

func (R *reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(R.data)
	R.data = R.data[4:]
	return i
}

func (R *reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(R.data)
	R.data = R.data[8:]
	return i
}

func (R *reader) readLuaInteger() int64 {
	return int64(R.readUint64())
}

func (R *reader) readLuaNumber() float64 {
	return math.Float64frombits(R.readUint64())
}

func (R *reader) readString() string {
	size := uint(R.readByte())
	if size == 0 {
		return ""
	}
	if size == 0xFF {
		size = uint(R.readUint64()) // size_t
	}
	bytes := R.readBytes(size - 1)
	return string(bytes) // todo
}

func (R *reader) checkHeader() {
	if string(R.readBytes(4)) != LUA_SIGNATURE {
		panic("not a precompiled chunk!")
	}
	if R.readByte() != LUAC_VERSION {
		panic("version mismatch!")
	}
	if R.readByte() != LUAC_FORMAT {
		panic("format mismatch!")
	}
	if string(R.readBytes(6)) != LUAC_DATA {
		panic("corrupted!")
	}
	if R.readByte() != CINT_SIZE {
		panic("int size mismatch!")
	}
	if R.readByte() != CSIZET_SIZE {
		panic("size_t size mismatch!")
	}
	if R.readByte() != INSTRUCTION_SIZE {
		panic("instruction size mismatch!")
	}
	if R.readByte() != LUA_INTEGER_SIZE {
		panic("lua_Integer size mismatch!")
	}
	if R.readByte() != LUA_NUMBER_SIZE {
		panic("lua_Number size mismatch!")
	}
	if R.readLuaInteger() != LUAC_INT {
		panic("endianness mismatch!")
	}
	if R.readLuaNumber() != LUAC_NUM {
		panic("float format mismatch!")
	}
}

func (R *reader) readProto(parentSource string) *Prototype {
	source := R.readString()
	if source == "" {
		source = parentSource
	}
	return &Prototype{
		Source:          source,
		LineDefined:     R.readUint32(),
		LastLineDefined: R.readUint32(),
		NumParams:       R.readByte(),
		IsVararg:        R.readByte(),
		MaxStackSize:    R.readByte(),
		Code:            R.readCode(),
		Constants:       R.readConstants(),
		Upvalues:        R.readUpvalues(),
		Protos:          R.readProtos(source),
		LineInfo:        R.readLineInfo(),
		LocVars:         R.readLocVars(),
		UpvalueNames:    R.readUpvalueNames(),
	}
}

func (R *reader) readCode() []uint32 {
	code := make([]uint32, R.readUint32())
	for i := range code {
		code[i] = R.readUint32()
	}
	return code
}

func (R *reader) readConstants() []interface{} {
	constants := make([]interface{}, R.readUint32())
	for i := range constants {
		constants[i] = R.readConstant()
	}
	return constants
}

func (R *reader) readConstant() interface{} {
	switch R.readByte() {
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return R.readByte() != 0
	case TAG_INTEGER:
		return R.readLuaInteger()
	case TAG_NUMBER:
		return R.readLuaNumber()
	case TAG_SHORT_STR, TAG_LONG_STR:
		return R.readString()
	default:
		panic("corrupted!") // todo
	}
}

func (R *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, R.readUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: R.readByte(),
			Idx:     R.readByte(),
		}
	}
	return upvalues
}

func (R *reader) readProtos(parentSource string) []*Prototype {
	protos := make([]*Prototype, R.readUint32())
	for i := range protos {
		protos[i] = R.readProto(parentSource)
	}
	return protos
}

func (R *reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32, R.readUint32())
	for i := range lineInfo {
		lineInfo[i] = R.readUint32()
	}
	return lineInfo
}

func (R *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, R.readUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: R.readString(),
			StartPC: R.readUint32(),
			EndPC:   R.readUint32(),
		}
	}
	return locVars
}

func (R *reader) readUpvalueNames() []string {
	names := make([]string, R.readUint32())
	for i := range names {
		names[i] = R.readString()
	}
	return names
}
