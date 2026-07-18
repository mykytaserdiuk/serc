package core

import (
	"unicode"

	"github.com/mykytaserdiuk/serc/ast"
)

type Lexer struct {
	source string
	cur    int
	row    int
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		source: source,
		cur:    0,
		row:    0,
	}
}

func (l *Lexer) isSpace() bool {
	// x:= l.source[l.cur]
	//	return rune(x) == '\n' || rune(x) == '\t' || rune(x) == '\r'
	return unicode.IsSpace(rune(l.source[l.cur]))
}

func (l *Lexer) isEmpty() bool {
	return !l.isNotEmpty()
}

func (l *Lexer) isNotEmpty() bool {
	return l.cur < len(l.source)
}

func (l *Lexer) moveRight() {
	for !l.isEmpty() && l.isSpace() {
		l.chop()
	}
}

func (l *Lexer) chop() {
	//fmt.Println(l.cur, len(l.source))
	if l.isNotEmpty() {
		x := l.source[l.cur]
		l.cur = l.cur + 1
		if x == '\n' {
			l.row += 1
		}
	}
}

func (l *Lexer) dropLine() {
	if l.isNotEmpty() && l.source[l.cur] != '\n' {
		l.chop()
	}
	if l.isEmpty() {
		l.chop()
	}
}

func (l *Lexer) NextToken() *ast.Token {
	l.moveRight()

	if l.isEmpty() {
		// end of file
		return &ast.Token{
			Value: "",
			Type_: ast.EOFTokenType,
			Line:  l.row,
		}
	}
	if l.source[l.cur] == ';' {
		for l.isNotEmpty() && l.source[l.cur] != '\n' {
			l.chop()
		}
		if l.isNotEmpty() {
			l.chop()
		}
		return l.NextToken()
	}

	first := l.source[l.cur]
	if unicode.IsLetter(rune(first)) {
		i := l.cur
		for l.isNotEmpty() && (l.source[l.cur] == '_' || unicode.IsLetter(rune(l.source[l.cur]))) { // || l.source[l.cur] != ' ') {
			l.chop()
		}
		letterTokens := map[string]ast.TokenType{
			"end":    ast.EndTokenType,
			"def":    ast.DefTokenType,
			"func":   ast.FuncTokenType,
			"if":     ast.IfTokenType,
			"else":   ast.ElseTokenType,
			"then":   ast.ThenTokenType,
			"return": ast.ReturnTokenType,
			"struct": ast.StructTokenType,
			"use":    ast.UseTokenType,
			"while":  ast.WhileTokenType,
		}
		value := l.source[i:l.cur]
		if val, ok := letterTokens[value]; ok {
			return &ast.Token{
				Type_: val,
				Value: value,
				Line:  l.row,
			}
		} else {
			return &ast.Token{
				Type_: ast.NameTokenType,
				Value: value,
				Line:  l.row,
			}
		}
	}

	if unicode.IsNumber(rune(first)) {
		i := l.cur
		for l.isNotEmpty() && unicode.IsNumber(rune(l.source[l.cur])) {
			l.chop()
		}
		return &ast.Token{
			Value: string(l.source[i:l.cur]),
			Type_: ast.NumTokenType,
		}
	}

	if first == '=' {
		next := l.source[l.cur+1]
		l.chop()
		if next == '=' {
			l.chop()
			return &ast.Token{
				Type_: ast.EqEqTokenType,
				Value: "==",
				Line:  l.row,
			}
		}
		return &ast.Token{
			Type_: ast.EqTokenType,
			Value: "=",
			Line:  l.row,
		}
	}

	unletterTokens := map[rune]ast.TokenType{
		'(': ast.OparenTokenType,
		')': ast.CparenTokenType,
		':': ast.ColonTokenType,
		',': ast.CommaTokenType,
		'>': ast.MoreTokenType,
		'<': ast.LessTokenType,
		'+': ast.PlusTokenType,
		'-': ast.MinusTokenType,
		'.': ast.DotTokenType,
	}
	if v, ok := unletterTokens[rune(first)]; ok {
		switch v {
		case ast.MoreTokenType:
			if l.source[l.cur+1] == '=' {
				l.chop()
				l.chop()
				return &ast.Token{
					Value: ">=",
					Type_: ast.EqMoreTokenType,
					Line:  l.row,
				}
			} else {
				l.chop()
				return &ast.Token{
					Value: ">",
					Type_: ast.MoreTokenType,
					Line:  l.row,
				}
			}
		case ast.LessTokenType:
			if l.source[l.cur+1] == '=' {
				l.chop()
				l.chop()
				return &ast.Token{
					Value: "<=",
					Type_: ast.EqLessTokenType,
					Line:  l.row,
				}
			} else if l.source[l.cur+1] == '>' {
				l.chop()
				l.chop()
				return &ast.Token{
					Value: "<>",
					Type_: ast.NotEqTokenType,
					Line:  l.row,
				}
			} else {
				l.chop()
				return &ast.Token{
					Value: "<",
					Type_: ast.LessTokenType,
					Line:  l.row,
				}
			}
		default:
			l.chop()
			return &ast.Token{
				Value: string(first),
				Type_: v,
				Line:  l.row,
			}
		}
	}

	if first == '"' {
		l.chop()
		value := []byte{}
		for l.isNotEmpty() && l.source[l.cur] != '"' {
			if l.source[l.cur] == '\\' {
				l.chop()
				if l.isEmpty() {
					panic("ERROR: unclosed \\ on the end of file")
				}

				switch l.source[l.cur] {
				case 'n':
					value = append(value, '\n')
				case 't':
					value = append(value, '\t')
				}
				l.chop()
			} else {
				value = append(value, l.source[l.cur])
				l.chop()
			}
		}

		if l.isNotEmpty() {
			//value = append(value, l.source[l.cur])
			l.chop()
			return &ast.Token{
				Type_: ast.StringTokenType,
				Value: string(value),
				Line:  l.row,
			}
		}
	}

	return &ast.Token{
		Value: "",
		Type_: ast.EOFTokenType,
		Line:  l.row,
	}
	//	panic("TODO next token")
}
