package core

type TokenType string

const (
	OparenTokenType TokenType = "oparen"
	CparenTokenType TokenType = "cparen"
	NameTokenType   TokenType = "name"
	ColonTokenType  TokenType = "colon"
	StringTokenType TokenType = "string"
	EndTokenType    TokenType = "end"
	FuncTokenType   TokenType = "func"
	DefTokenType    TokenType = "def"
	NumTokenType    TokenType = "num"
	EqTokenType     TokenType = "eq"
	CommaTokenType  TokenType = "comma"
	DotTokenType    TokenType = "dot"

	IfTokenType   TokenType = "if"
	ThenTokenType TokenType = "then"
	ElseTokenType TokenType = "else"
	WhileTokenType TokenType = "while"

	StructTokenType TokenType = "struct"
	UseTokenType    TokenType = "use"

	MoreTokenType   TokenType = "more"
	LessTokenType   TokenType = "less"
	EqLessTokenType TokenType = "eqless"
	EqMoreTokenType TokenType = "eqmore"
	EqEqTokenType   TokenType = "eqeq"
	NotEqTokenType  TokenType = "noteq"

	PlusTokenType  TokenType = "plus"
	MinusTokenType TokenType = "minus"

	ReturnTokenType TokenType = "return"
	EOFTokenType    TokenType = "EOF"
)

type Token struct {
	type_ TokenType
	value string
	line  int
}
