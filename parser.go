package main

import (
	"fmt"
)

type Node interface{}
type Statement interface{
	Node
}

type Expression interface{
	Node
}

type Func struct{
	name string
	params []string
	body []Statement
}

type FuncCall struct{
	name string
	args []string
}
func (f FuncCall) node(){}
func (f FuncCall) statement(){}


type Parser struct{
	l *Lexer
}

func (p *Parser) expectType(expectedTypes ...TokenType) (*Token, bool) {
	token, ok := p.l.NextToken()
	if ok {
		for _, t := range expectedTypes{
			if token.type_ ==  t {
				return token, true
			}
		}
	}
	fmt.Println("Expected: ", expectedTypes, " got: ", token.type_)
	return nil, false
}

func (p *Parser) parseType() string {
	returnType, ok := p.expectType(NameTokenType)
	if !ok {
		 fmt.Printf("%+v: ERROR: unexpected type %s",
			 returnType, returnType.value);
		return  ""
	}
	return string(returnType.type_)
}

func (p *Parser) parseFunc() *Func{
	token, ok := p.expectType(FuncTokenType)
	if !ok {
		fmt.Printf("ERROR: parse func: expected type: %s, got: %s", FuncTokenType, token.type_)
		return nil
	}
	nameToken, ok := p.expectType(NameTokenType)
	if !ok {
		fmt.Printf("ERROR: parse func: expected type: %s, got: %s", NameTokenType, nameToken.type_)
		return nil
	}

	params := p.parseParams()
	block := p.parseBlock()

	return &Func{
		name: nameToken.value,
		body: block,
		params: params,
	}
}

func (p *Parser) parseBlock() ([]Statement) {
	if _, ok := p.expectType(ColonTokenType) ; !ok {
		return nil
	}

	block := make([]Statement, 0)
	for {
		token, ok := p.expectType(NameTokenType, EndTokenType)
		if !ok {
			fmt.Printf("ERROR: parse block: expected types: %s, %s", NameTokenType, EndTokenType)
			break
		}
		end := false
		switch token.type_{
			case EndTokenType:
			end = true;
			case NameTokenType:
			params := p.parseParams()
			block = append(block, FuncCall{
				name:token.value,
				args: params,
			})
		}
		if end {
			break
		}
	}
	return block
}

func (p *Parser) parseParams() []string {
	_, ok := p.expectType(OparenTokenType)
	params := make([]string, 0)
	if !ok {
		return params
	}
	for {
		end := false
		if token, ok := p.expectType(CparenTokenType, StringTokenType) ; ok{
			switch token.type_{
				case CparenTokenType:
				 end = true
				case StringTokenType:
				params = append(params, token.value)
			}
		}
		if end {
			break
		}
	}
	return params
}
