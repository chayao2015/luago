package ast

// chunk ::= block
// type Chunk *Block

// block ::= {stat} [retstat]
// retstat ::= return [explist] [‘;’]
// explist ::= exp {‘,’ exp}

// 在EBNF里，“::=”表示“被定义为”的意思 {A}表示A可以出现任意次（0次或多次），[A]表示A可选（可以出现0次或1次）

// 代码块是任意多条语句再加上一条可选的返回语句；
// 返回语句是关键字return后跟可选的表达式列表，以及一个可选的分号；表达式列表则是1到多个表达式，由逗号分隔

type Block struct {
	LastLine int //记录代码块的末尾行号
	Stats    []Stat
	RetExps  []Exp
}

/*
	Lua语法EBNF描述

chunk ::= block

block ::= {stat} [retstat]

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

retstat ::= return [explist] [‘;’]

label ::= ‘::’ Name ‘::’

funcname ::= Name {‘.’ Name} [‘:’ Name]

varlist ::= var {‘,’ var}

var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name

namelist ::= Name {‘,’ Name}

explist ::= exp {‘,’ exp}

exp ::=  nil | false | true | Numeral | LiteralString | ‘...’ | functiondef |
	prefixexp | tableconstructor | exp binop exp | unop exp

prefixexp ::= var | functioncall | ‘(’ exp ‘)’

functioncall ::=  prefixexp args | prefixexp ‘:’ Name args

args ::=  ‘(’ [explist] ‘)’ | tableconstructor | LiteralString

functiondef ::= function funcbody

funcbody ::= ‘(’ [parlist] ‘)’ block end

parlist ::= namelist [‘,’ ‘...’] | ‘...’

tableconstructor ::= ‘{’ [fieldlist] ‘}’

fieldlist ::= field {fieldsep field} [fieldsep]

field ::= ‘[’ exp ‘]’ ‘=’ exp | Name ‘=’ exp | exp

fieldsep ::= ‘,’ | ‘;’

binop ::=  ‘+’ | ‘-’ | ‘*’ | ‘/’ | ‘//’ | ‘^’ | ‘%’ |
	‘&’ | ‘~’ | ‘|’ | ‘>>’ | ‘<<’ | ‘..’ |
	‘<’ | ‘<=’ | ‘>’ | ‘>=’ | ‘==’ | ‘~=’ |
	and | or

unop ::= ‘-’ | not | ‘#’ | ‘~’
*/
