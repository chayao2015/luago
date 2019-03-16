package ast

/*
exp ::=  nil | false | true | Numeral | LiteralString | ‘...’ | functiondef |
	prefixexp | tableconstructor | exp binop exp | unop exp

prefixexp ::= var | functioncall | ‘(’ exp ‘)’

var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name

functioncall ::=  prefixexp args | prefixexp ‘:’ Name args
*/

type Exp interface{}

type NilExp struct{ Line int }
type TrueExp struct{ Line int }
type FalseExp struct{ Line int }
type VarargExp struct{ Line int }

type IntegerExp struct {
	Line int
	Val  int64
}

type FloatExp struct {
	Line int
	Val  float64
}

// LiteralString
type StringExp struct {
	Line int
	Str  string
}

type NameExp struct {
	Line int
	Name string
}

// prefixexp ::= var | functioncall | ‘(’ exp ‘)’
// functioncall ::=  prefixexp args | prefixexp ‘:’ Name args
type FuncCallExp struct {
	Line      int // line of `(` ?
	LastLine  int // line of ')'
	PrefixExp Exp
	NameExp   *StringExp
	Args      []Exp
}

// functiondef ::= function funcbody
// funcbody ::= ‘(’ [parlist] ‘)’ block end
// parlist ::= namelist [‘,’ ‘...’] | ‘...’
// namelist ::= Name {‘,’ Name}
type FuncDefExp struct {
	Line     int // line of `local funciton`
	LastLine int // line of 'end'
	ParList  []string
	IsVararg bool
	Block    *Block
}

type UnopExp struct {
	Line int // line of operator
	Op   int // operator
	Exp  Exp
}

type BinopExp struct {
	Line int // line of operator
	Op   int // operator
	Exp1 Exp
	Exp2 Exp
}

type ConcatExp struct {
	Line int // line of last ..
	Exps []Exp
}

// tableconstructor ::= ‘{’ [fieldlist] ‘}’

// fieldlist ::= field {fieldsep field} [fieldsep]

// field ::= ‘[’ exp ‘]’ ‘=’ exp | Name ‘=’ exp | exp

// fieldsep ::= ‘,’ | ‘;’
type TableConstructorExp struct {
	Line     int
	LastLine int
	KeyExps  []Exp
	ValExps  []Exp
}

// prefixexp ::= var | functioncall | ‘(’ exp ‘)’
//var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name
// functioncall ::=  prefixexp args | prefixexp ‘:’ Name args
// \\\\\
// prefixexp ::= Name
// | '(' exp ')'
// | prefixexp '[' exp ']'
// | prefixexp '.' Name
// | prefixexp [':' Name] args

type ParensExp struct {
	Exp Exp
}

type TableAccessExp struct {
	LastLine  int // line of `]` ?
	PrefixExp Exp
	KeyExp    Exp
}
