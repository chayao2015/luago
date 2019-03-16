package ast

// 在命令式编程语言里，语句（Statement）是最基本的执行单位，表达式
// （Expression）则是构成语句的要素之一。语句和表达式的主要区别在于：语句只能
// 执行不能用于求值，而表达式只能用于求值不能单独执行
// 语句和表达式也并非泾渭分明，比如在Lua里，函数调用既可以是表达式，也可以是语句

/*
stat ::=  ‘;’ |
	varlist ‘=’ explist |
	functioncall |
	label |
	break |
	goto Name |
	do block end |
	while exp do block end |
	repeat block until exp |
	if exp then block {elseif exp then block} [else block] end |
	for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end |
	for namelist in explist do block end |
	function funcname funcbody |
	local function Name funcbody |
	local namelist [‘=’ explist]
*/

type Stat interface{}

type EmptyStat struct{}              // ‘;’ 空语句
type BreakStat struct{ Line int }    // break
type LabelStat struct{ Name string } // ‘::’ Name ‘::’ 标签
type GotoStat struct{ Name string }  // goto Name
type DoStat struct{ Block *Block }   // do block end
type FuncCallStat = FuncCallExp      // functioncall

// while exp do block end
type WhileStat struct {
	Exp   Exp
	Block *Block
}

type RepeatStat struct {
	Block *Block
	Exp   Exp
}

// 我们把表达式收集到Exps字段里，把语句块收集到Blocks字段里。表达式和
// 语句块按索引一一对应，索引0处是if-then表达式和块，其余索引处是elseif-then表达式和块
// if exp then block {elseif exp then block} [else block] end
type IfStat struct {
	Exps   []Exp
	Blocks []*Block
}

// for Name '=’ exp ',’ exp [',’ exp] do block end
//需要把关键字for和do所在的行号记录下来，以供代码生成阶段使用
type ForNumStat struct {
	LineOfFor int
	LineOfDo  int
	VarName   string
	InitExp   Exp
	LimitExp  Exp
	StepExp   Exp
	Block     *Block
}

// for namelist in explist do block end
// namelist ::= Name {',' Name}
// explist ::= exp {',' exp
type ForInStat struct {
	LineOfDo int
	NameList []string
	ExpList  []Exp
	Block    *Block
}

// varlist ‘=’ explist
// varlist ::= var {‘,’ var}
// var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name
type AssignStat struct {
	LastLine int
	VarList  []Exp
	ExpList  []Exp
}

// local namelist [‘=’ explist]
// namelist ::= Name {‘,’ Name}
// explist ::= exp {‘,’ exp}
type LocalVarDeclStat struct {
	LastLine int
	NameList []string
	ExpList  []Exp
}

// local function Name funcbody
type LocalFuncDefStat struct {
	Name string
	Exp  *FuncDefExp
}
