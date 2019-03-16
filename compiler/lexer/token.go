package lexer

/*Lua词法规则(token)*/
// token按其作用可以分为不同的类型，比较常见的类型注释、关键字、标识符、字面量、运算符、分隔符等

// 标识符（Identifier）主要用来命名变量。Lua标识符以字母或下划线开头，后跟
// 数字、字母或者下划线的任意组合。Lua是大小写敏感语言，因此var、Var和VAR是
// 三个不同的标识符。按照惯例，应该避免在程序中使用以下划线开头，后跟一个
// 或多个大写字母的标识符（比如_ENV）

// 关键字（Keyword）具有特殊含义，由编程语言保留，不能当作标识符使用。下
// 面是Lua语言所保留的关键字（共22个
// and break do else elseif end
// false for function goto if in
// local nil not or repeat return
// then true until while

// 数字字面量
// 最简单的是十进制整数写法，比如3、314。当使用小数写法时，整数部分和小数部分都可以省略，比如3.、3.14、.14。
// 还可以加上指数部分写成科学计数法，比如0.314E1、314e-2

// 十六进制写法以0x或者0X开头，比如0xff、0X3A、0x3.243F6A8885A。十六进
// 制也可以使用科学计数法，但是指数部分用字母p（或者P）表示，只能使用十进制数字，并且表示的是2的多少次方，比如0xA23p-4

// 如果数字字面量不包含小数和指数部分，也没有超出Lua整数的表示范围，则
// 会被Lua解释成整数值，否则会被Lua解释成浮点数值

// 字符串字面量
// Lua字符串字面量分为长字符串和短字符串两种  短字符串"" ''里面可以包含转义序列 长字符串[[]]

//运算符和分隔符

// token kind
const (
	TOKEN_EOF         = iota           // end-of-file
	TOKEN_VARARG                       // ...
	TOKEN_SEP_SEMI                     // ;
	TOKEN_SEP_COMMA                    // ,
	TOKEN_SEP_DOT                      // .
	TOKEN_SEP_COLON                    // :
	TOKEN_SEP_LABEL                    // ::
	TOKEN_SEP_LPAREN                   // (
	TOKEN_SEP_RPAREN                   // )
	TOKEN_SEP_LBRACK                   // [
	TOKEN_SEP_RBRACK                   // ]
	TOKEN_SEP_LCURLY                   // {
	TOKEN_SEP_RCURLY                   // }
	TOKEN_OP_ASSIGN                    // =
	TOKEN_OP_MINUS                     // - (sub or unm)
	TOKEN_OP_WAVE                      // ~ (bnot or bxor)
	TOKEN_OP_ADD                       // +
	TOKEN_OP_MUL                       // *
	TOKEN_OP_DIV                       // /
	TOKEN_OP_IDIV                      // //
	TOKEN_OP_POW                       // ^
	TOKEN_OP_MOD                       // %
	TOKEN_OP_BAND                      // &
	TOKEN_OP_BOR                       // |
	TOKEN_OP_SHR                       // >>
	TOKEN_OP_SHL                       // <<
	TOKEN_OP_CONCAT                    // ..
	TOKEN_OP_LT                        // <
	TOKEN_OP_LE                        // <=
	TOKEN_OP_GT                        // >
	TOKEN_OP_GE                        // >=
	TOKEN_OP_EQ                        // ==
	TOKEN_OP_NE                        // ~=
	TOKEN_OP_LEN                       // #
	TOKEN_OP_AND                       // and
	TOKEN_OP_OR                        // or
	TOKEN_OP_NOT                       // not
	TOKEN_KW_BREAK                     // break
	TOKEN_KW_DO                        // do
	TOKEN_KW_ELSE                      // else
	TOKEN_KW_ELSEIF                    // elseif
	TOKEN_KW_END                       // end
	TOKEN_KW_FALSE                     // false
	TOKEN_KW_FOR                       // for
	TOKEN_KW_FUNCTION                  // function
	TOKEN_KW_GOTO                      // goto
	TOKEN_KW_IF                        // if
	TOKEN_KW_IN                        // in
	TOKEN_KW_LOCAL                     // local
	TOKEN_KW_NIL                       // nil
	TOKEN_KW_REPEAT                    // repeat
	TOKEN_KW_RETURN                    // return
	TOKEN_KW_THEN                      // then
	TOKEN_KW_TRUE                      // true
	TOKEN_KW_UNTIL                     // until
	TOKEN_KW_WHILE                     // while
	TOKEN_IDENTIFIER                   // identifier
	TOKEN_NUMBER                       // number literal
	TOKEN_STRING                       // string literal
	TOKEN_OP_UNM      = TOKEN_OP_MINUS // unary minus
	TOKEN_OP_SUB      = TOKEN_OP_MINUS
	TOKEN_OP_BNOT     = TOKEN_OP_WAVE
	TOKEN_OP_BXOR     = TOKEN_OP_WAVE
)

var keywords = map[string]int{
	"and":      TOKEN_OP_AND,
	"break":    TOKEN_KW_BREAK,
	"do":       TOKEN_KW_DO,
	"else":     TOKEN_KW_ELSE,
	"elseif":   TOKEN_KW_ELSEIF,
	"end":      TOKEN_KW_END,
	"false":    TOKEN_KW_FALSE,
	"for":      TOKEN_KW_FOR,
	"function": TOKEN_KW_FUNCTION,
	"goto":     TOKEN_KW_GOTO,
	"if":       TOKEN_KW_IF,
	"in":       TOKEN_KW_IN,
	"local":    TOKEN_KW_LOCAL,
	"nil":      TOKEN_KW_NIL,
	"not":      TOKEN_OP_NOT,
	"or":       TOKEN_OP_OR,
	"repeat":   TOKEN_KW_REPEAT,
	"return":   TOKEN_KW_RETURN,
	"then":     TOKEN_KW_THEN,
	"true":     TOKEN_KW_TRUE,
	"until":    TOKEN_KW_UNTIL,
	"while":    TOKEN_KW_WHILE,
}
