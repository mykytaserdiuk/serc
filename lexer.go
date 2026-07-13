package main

import (
	"unicode"
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

func (l *Lexer) NextToken() (*Token) {
	l.moveRight()

	if l.isEmpty() {
		// end of file
		return &Token{
			value: "",
			type_: EOFTokenType,
		}
	}
	if l.source[l.cur] == ';' {
		for l.isNotEmpty() && l.source[l.cur] != '\n'{
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
		for l.isNotEmpty() && unicode.IsLetter(rune(l.source[l.cur])) {
			l.chop()
		}
		letterTokens := map[string]TokenType{
			"end":  EndTokenType,
			"func": FuncTokenType,
			"def":  DefTokenType,
			"if": IfTokenType,
			"else": ElseTokenType,
			"then":ThenTokenType,
			"return": ReturnTokenType,
		}
		value := l.source[i:l.cur]
		if val, ok := letterTokens[value]; ok {
			return &Token{
				type_: val,
				value: value,
			}
		} else {
			return &Token{
				type_: NameTokenType,
				value: value,
			}
		}
	}

	if unicode.IsNumber(rune(first)) {
		i := l.cur
		for l.isNotEmpty() && unicode.IsNumber(rune(l.source[l.cur])) {
			l.chop()
		}
		return &Token{
			value: string(l.source[i:l.cur]),
			type_: NumTokenType,
		}
	}

	if first == '=' {
		l.chop()
		return &Token{
			type_: EqTokenType,
			value: "=",
		}
	}

	unletterTokens := map[rune]TokenType{
		'(': OparenTokenType,
		')': CparenTokenType,
		':': ColonTokenType,
		',': CommaTokenType,
		'>':MoreTokenType,
		'<':LessTokenType,
	}
	if v, ok := unletterTokens[rune(first)]; ok {
		l.chop()
		switch v{
			case MoreTokenType:
			if l.source[l.cur+1] == '=' {
				l.chop()
				return &Token{
					value: ">=",
					type_: EqMoreTokenType,
				}
			} else {
				return &Token{
					value: ">",
					type_: MoreTokenType,
				}
			}
			case LessTokenType:
			if l.source[l.cur+1] == '=' {
				l.chop()
				return &Token{
					value: "<=",
					type_: EqLessTokenType,
				}
			} else {
				return &Token{
					value: "<",
					type_: LessTokenType,
				}
			}
			default:
			return &Token{
				value: string(first),
				type_: v,
			}
		}
	}

	if first == '"' {
		l.chop()
		value := []byte{}
		for l.isNotEmpty() && l.source[l.cur] != '"' {
			if l.source[l.cur] == '\\'{
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
			} else{
				value = append(value, l.source[l.cur])
				l.chop()
			}
		}

		if l.isNotEmpty() {
			//value = append(value, l.source[l.cur])
			l.chop()
			return &Token{
				type_: StringTokenType,
				value: string(value),
			}
		}
	}

	return &Token{
		value: "",
		type_: EOFTokenType,
	}
	//	panic("TODO next token")
}
