package main

import (
	"unicode"
)


type Lexer struct{
	source string
	cur int
	row int
}

func NewLexer(source string) *Lexer{
	return &Lexer{
		source: source,
		cur: 0,
		row: 0,
	}
}

func (l *Lexer) isSpace() bool{
	// x:= l.source[l.cur]
	//	return rune(x) == '\n' || rune(x) == '\t' || rune(x) == '\r'
	return unicode.IsSpace(rune(l.source[l.cur]))
}

func (l *Lexer) isEmpty() bool{
	return !l.isNotEmpty()
}

func (l *Lexer) isNotEmpty() bool{
	return l.cur < len(l.source)
}

func (l *Lexer) moveRight(){
	for(!l.isEmpty() && l.isSpace()) {
		l.chop()
	}
}

func (l *Lexer) chop(){
	//fmt.Println(l.cur, len(l.source))
	if(l.isNotEmpty()){
		x:= l.source[l.cur]
		l.cur = l.cur +1
		if (x == '\n'){
			l.row+=1
		}
	}
}

func (l *Lexer)  dropLine() {
	if (l.isNotEmpty() && l.source[l.cur] != '\n'){
		l.chop()
	}
	if (l.isEmpty()) {
		l.chop()
	}
}

func (l *Lexer) NextToken() (*Token, bool){
	l.moveRight()

	if l.isEmpty() {
		// end of file
		return &Token{}, false
	}

	first := l.source[l.cur]
	//fmt.Print(first)
	if unicode.IsLetter(rune(first)){
		i := l.cur
		for l.isNotEmpty() && unicode.IsLetter(rune(l.source[l.cur])){
			l.chop()
		}
		letterTokens:= map[string]TokenType{
			"end":EndTokenType,
			"func":FuncTokenType,
		}
		value := l.source[i:l.cur]
		if val, ok := letterTokens[value]; ok {
			return &Token{
				type_: val,
				value: value,
			}, true
		} else {
			return &Token{
				type_ : NameTokenType,
				value: value,
			}, true
		}
	}

	unletterTokens:= map[rune]TokenType{
		'(':OparenTokenType,
		')':CparenTokenType,
		':':ColonTokenType,
	}
	if v, ok := unletterTokens[rune(first)]; ok {
		l.chop()
		return &Token{
			value: string(first),
			type_: v,
		}, true
	}

	if first == '"'  {
		l.chop()
		i:= l.cur
		for l.isNotEmpty() && l.source[l.cur] != '"'{
			l.chop()
		}
		if l.isNotEmpty() {
			value := l.source[i:l.cur]
			l.chop()
			return &Token{
				type_: StringTokenType,
				value: value,
			}, true
		}
	}

	return  nil, false
	//	panic("TODO next token")
}
