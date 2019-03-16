package lexer

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// 正则表达式来处理换行符序列
// var reNewLine = regexp.MustCompile("\r\n|\n\r|\n|\r")
// var reIdentifier = regexp.MustCompile(`^[_\d\w]+`)
// var reNumber = regexp.MustCompile(`^0[xX][0-9a-fA-F]*(\.[0-9a-fA-F]*)?([pP][+\-]?[0-9]+)?|^[0-9]*(\.[0-9]*)?([eE][+\-]?[0-9]+)?`)
// var reShortStr = regexp.MustCompile(`(?s)(^'(\\\\|\\'|\\\n|\\z\s*|[^'\n])*')|(^"(\\\\|\\"|\\\n|\\z\s*|[^"\n])*")`)

// //正则表达式来探测左长方括号
// var reOpeningLongBracket = regexp.MustCompile(`^\[=*\[`)

// //10进制数字正则
// var reDecEscapeSeq = regexp.MustCompile(`^\\[0-9]{1,3}`)

// //16进制数字正则
// var reHexEscapeSeq = regexp.MustCompile(`^\\x[0-9a-fA-F]{2}`)

// //Unicode字符正则
// var reUnicodeEscapeSeq = regexp.MustCompile(`^\\u\{[0-9a-fA-F]+\}`)

/*词法分析器*/

// 编译器在编译源代码时，也不是以字符为单位，而是以“token”为
// 单位进行处理的。词法分析器的作用就是根据编程语言的词法规则，把源代码（字符流）分解为token流

// chunk line 构成词法分析内部状态
// chunkNmae 仅用于在词法分析过程出错时生成错误信息

var reNewLine = regexp.MustCompile("\r\n|\n\r|\n|\r")
var reIdentifier = regexp.MustCompile(`^[_\d\w]+`)
var reNumber = regexp.MustCompile(`^0[xX][0-9a-fA-F]*(\.[0-9a-fA-F]*)?([pP][+\-]?[0-9]+)?|^[0-9]*(\.[0-9]*)?([eE][+\-]?[0-9]+)?`)
var reShortStr = regexp.MustCompile(`(?s)(^'(\\\\|\\'|\\\n|\\z\s*|[^'\n])*')|(^"(\\\\|\\"|\\\n|\\z\s*|[^"\n])*")`)
var reOpeningLongBracket = regexp.MustCompile(`^\[=*\[`)

var reDecEscapeSeq = regexp.MustCompile(`^\\[0-9]{1,3}`)
var reHexEscapeSeq = regexp.MustCompile(`^\\x[0-9a-fA-F]{2}`)
var reUnicodeEscapeSeq = regexp.MustCompile(`^\\u\{[0-9a-fA-F]+\}`)

type Lexer struct {
	chunk         string // source code
	chunkName     string // source name
	line          int    // current line number
	nextToken     string
	nextTokenKind int
	nextTokenLine int
}

func NewLexer(chunk, chunkName string) *Lexer {
	return &Lexer{chunk, chunkName, 1, "", 0, 0}
}

func (L *Lexer) Line() int {
	return L.line
}

func (L *Lexer) LookAhead() int {
	if L.nextTokenLine > 0 {
		return L.nextTokenKind
	}
	currentLine := L.line
	line, kind, token := L.NextToken()
	L.line = currentLine
	L.nextTokenLine = line
	L.nextTokenKind = kind
	L.nextToken = token
	return kind
}

func (L *Lexer) NextIdentifier() (line int, token string) {
	return L.NextTokenOfKind(TOKEN_IDENTIFIER)
}

func (L *Lexer) NextTokenOfKind(kind int) (line int, token string) {
	line, _kind, token := L.NextToken()
	if kind != _kind {
		L.error("syntax error near '%s'", token)
	}
	return line, token
}

func (L *Lexer) NextToken() (line, kind int, token string) {
	if L.nextTokenLine > 0 {
		line = L.nextTokenLine
		kind = L.nextTokenKind
		token = L.nextToken
		L.line = L.nextTokenLine
		L.nextTokenLine = 0
		return
	}

	L.skipWhiteSpaces()
	if len(L.chunk) == 0 {
		return L.line, TOKEN_EOF, "EOF"
	}

	switch L.chunk[0] {
	case ';':
		L.next(1)
		return L.line, TOKEN_SEP_SEMI, ";"
	case ',':
		L.next(1)
		return L.line, TOKEN_SEP_COMMA, ","
	case '(':
		L.next(1)
		return L.line, TOKEN_SEP_LPAREN, "("
	case ')':
		L.next(1)
		return L.line, TOKEN_SEP_RPAREN, ")"
	case ']':
		L.next(1)
		return L.line, TOKEN_SEP_RBRACK, "]"
	case '{':
		L.next(1)
		return L.line, TOKEN_SEP_LCURLY, "{"
	case '}':
		L.next(1)
		return L.line, TOKEN_SEP_RCURLY, "}"
	case '+':
		L.next(1)
		return L.line, TOKEN_OP_ADD, "+"
	case '-':
		L.next(1)
		return L.line, TOKEN_OP_MINUS, "-"
	case '*':
		L.next(1)
		return L.line, TOKEN_OP_MUL, "*"
	case '^':
		L.next(1)
		return L.line, TOKEN_OP_POW, "^"
	case '%':
		L.next(1)
		return L.line, TOKEN_OP_MOD, "%"
	case '&':
		L.next(1)
		return L.line, TOKEN_OP_BAND, "&"
	case '|':
		L.next(1)
		return L.line, TOKEN_OP_BOR, "|"
	case '#':
		L.next(1)
		return L.line, TOKEN_OP_LEN, "#"
	case ':':
		if L.test("::") {
			L.next(2)
			return L.line, TOKEN_SEP_LABEL, "::"
		} else {
			L.next(1)
			return L.line, TOKEN_SEP_COLON, ":"
		}
	case '/':
		if L.test("//") {
			L.next(2)
			return L.line, TOKEN_OP_IDIV, "//"
		} else {
			L.next(1)
			return L.line, TOKEN_OP_DIV, "/"
		}
	case '~':
		if L.test("~=") {
			L.next(2)
			return L.line, TOKEN_OP_NE, "~="
		} else {
			L.next(1)
			return L.line, TOKEN_OP_WAVE, "~"
		}
	case '=':
		if L.test("==") {
			L.next(2)
			return L.line, TOKEN_OP_EQ, "=="
		} else {
			L.next(1)
			return L.line, TOKEN_OP_ASSIGN, "="
		}
	case '<':
		if L.test("<<") {
			L.next(2)
			return L.line, TOKEN_OP_SHL, "<<"
		} else if L.test("<=") {
			L.next(2)
			return L.line, TOKEN_OP_LE, "<="
		} else {
			L.next(1)
			return L.line, TOKEN_OP_LT, "<"
		}
	case '>':
		if L.test(">>") {
			L.next(2)
			return L.line, TOKEN_OP_SHR, ">>"
		} else if L.test(">=") {
			L.next(2)
			return L.line, TOKEN_OP_GE, ">="
		} else {
			L.next(1)
			return L.line, TOKEN_OP_GT, ">"
		}
	case '.':
		if L.test("...") {
			L.next(3)
			return L.line, TOKEN_VARARG, "..."
		} else if L.test("..") {
			L.next(2)
			return L.line, TOKEN_OP_CONCAT, ".."
		} else if len(L.chunk) == 1 || !isDigit(L.chunk[1]) {
			L.next(1)
			return L.line, TOKEN_SEP_DOT, "."
		}
	case '[':
		if L.test("[[") || L.test("[=") {
			return L.line, TOKEN_STRING, L.scanLongString()
		} else {
			L.next(1)
			return L.line, TOKEN_SEP_LBRACK, "["
		}
	case '\'', '"':
		return L.line, TOKEN_STRING, L.scanShortString()
	}

	c := L.chunk[0]
	if c == '.' || isDigit(c) {
		token := L.scanNumber()
		return L.line, TOKEN_NUMBER, token
	}
	if c == '_' || isLetter(c) {
		token := L.scanIdentifier()
		if kind, found := keywords[token]; found {
			return L.line, kind, token // keyword
		} else {
			return L.line, TOKEN_IDENTIFIER, token
		}
	}

	L.error("unexpected symbol near %q", c)
	return
}

func (L *Lexer) next(n int) {
	L.chunk = L.chunk[n:]
}

func (L *Lexer) test(s string) bool {
	return strings.HasPrefix(L.chunk, s)
}

func (L *Lexer) error(f string, a ...interface{}) {
	err := fmt.Sprintf(f, a...)
	err = fmt.Sprintf("%s:%d: %s", L.chunkName, L.line, err)
	panic(err)
}

func (L *Lexer) skipWhiteSpaces() {
	for len(L.chunk) > 0 {
		if L.test("--") {
			L.skipComment()
		} else if L.test("\r\n") || L.test("\n\r") {
			L.next(2)
			L.line += 1
		} else if isNewLine(L.chunk[0]) {
			L.next(1)
			L.line += 1
		} else if isWhiteSpace(L.chunk[0]) {
			L.next(1)
		} else {
			break
		}
	}
}

func (L *Lexer) skipComment() {
	L.next(2) // skip --

	// long comment ?
	if L.test("[") {
		if reOpeningLongBracket.FindString(L.chunk) != "" {
			L.scanLongString()
			return
		}
	}

	// short comment
	for len(L.chunk) > 0 && !isNewLine(L.chunk[0]) {
		L.next(1)
	}
}

func (L *Lexer) scanIdentifier() string {
	return L.scan(reIdentifier)
}

func (L *Lexer) scanNumber() string {
	return L.scan(reNumber)
}

func (L *Lexer) scan(re *regexp.Regexp) string {
	if token := re.FindString(L.chunk); token != "" {
		L.next(len(token))
		return token
	}
	panic("unreachable!")
}

func (L *Lexer) scanLongString() string {
	openingLongBracket := reOpeningLongBracket.FindString(L.chunk)
	if openingLongBracket == "" {
		L.error("invalid long string delimiter near '%s'",
			L.chunk[0:2])
	}

	closingLongBracket := strings.Replace(openingLongBracket, "[", "]", -1)
	closingLongBracketIdx := strings.Index(L.chunk, closingLongBracket)
	if closingLongBracketIdx < 0 {
		L.error("unfinished long string or comment")
	}

	str := L.chunk[len(openingLongBracket):closingLongBracketIdx]
	L.next(closingLongBracketIdx + len(closingLongBracket))

	str = reNewLine.ReplaceAllString(str, "\n")
	L.line += strings.Count(str, "\n")
	if len(str) > 0 && str[0] == '\n' {
		str = str[1:]
	}

	return str
}

func (L *Lexer) scanShortString() string {
	if str := reShortStr.FindString(L.chunk); str != "" {
		L.next(len(str))
		str = str[1 : len(str)-1]
		if strings.Index(str, `\`) >= 0 {
			L.line += len(reNewLine.FindAllString(str, -1))
			str = L.escape(str)
		}
		return str
	}
	L.error("unfinished string")
	return ""
}

func (L *Lexer) escape(str string) string {
	var buf bytes.Buffer

	for len(str) > 0 {
		if str[0] != '\\' {
			buf.WriteByte(str[0])
			str = str[1:]
			continue
		}

		if len(str) == 1 {
			L.error("unfinished string")
		}

		switch str[1] {
		case 'a':
			buf.WriteByte('\a')
			str = str[2:]
			continue
		case 'b':
			buf.WriteByte('\b')
			str = str[2:]
			continue
		case 'f':
			buf.WriteByte('\f')
			str = str[2:]
			continue
		case 'n', '\n':
			buf.WriteByte('\n')
			str = str[2:]
			continue
		case 'r':
			buf.WriteByte('\r')
			str = str[2:]
			continue
		case 't':
			buf.WriteByte('\t')
			str = str[2:]
			continue
		case 'v':
			buf.WriteByte('\v')
			str = str[2:]
			continue
		case '"':
			buf.WriteByte('"')
			str = str[2:]
			continue
		case '\'':
			buf.WriteByte('\'')
			str = str[2:]
			continue
		case '\\':
			buf.WriteByte('\\')
			str = str[2:]
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // \ddd
			if found := reDecEscapeSeq.FindString(str); found != "" {
				d, _ := strconv.ParseInt(found[1:], 10, 32)
				if d <= 0xFF {
					buf.WriteByte(byte(d))
					str = str[len(found):]
					continue
				}
				L.error("decimal escape too large near '%s'", found)
			}
		case 'x': // \xXX
			if found := reHexEscapeSeq.FindString(str); found != "" {
				d, _ := strconv.ParseInt(found[2:], 16, 32)
				buf.WriteByte(byte(d))
				str = str[len(found):]
				continue
			}
		case 'u': // \u{XXX}
			if found := reUnicodeEscapeSeq.FindString(str); found != "" {
				d, err := strconv.ParseInt(found[3:len(found)-1], 16, 32)
				if err == nil && d <= 0x10FFFF {
					buf.WriteRune(rune(d))
					str = str[len(found):]
					continue
				}
				L.error("UTF-8 value too large near '%s'", found)
			}
		case 'z':
			str = str[2:]
			for len(str) > 0 && isWhiteSpace(str[0]) { // todo
				str = str[1:]
			}
			continue
		}
		L.error("invalid escape sequence near '\\%c'", str[1])
	}

	return buf.String()
}

func isWhiteSpace(c byte) bool {
	switch c {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	}
	return false
}

func isNewLine(c byte) bool {
	return c == '\r' || c == '\n'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isLetter(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}
