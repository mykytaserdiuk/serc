package main

type TokenType string
const (
	OparenTokenType TokenType = "oparen"
	CparenTokenType TokenType = "cparen"
	NameTokenType TokenType = "name"
	ColonTokenType TokenType = "colon"
	StringTokenType TokenType = "string"
	EndTokenType TokenType = "end"
	FuncTokenType TokenType = "func"
)

type Token struct{
	type_ TokenType
	value string
}
