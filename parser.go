package main

import (
	"fmt"
	"strconv"
)

type Node interface{}
type Statement interface {
	Node
	statement()
}

type Expression interface {
	Node
	expression()
}

type Func struct {
	name   string
	params []string
	body   []Statement
}

type Parser struct {
	l *Lexer
}

func (p *Parser) expectType(expectedTypes ...TokenType) (*Token, bool) {
	token, ok := p.l.NextToken()
	if ok {
		for _, t := range expectedTypes {
			if token.type_ == t {
				return token, true
			}
		}
	}
	return token, false
}

func (p *Parser) parseType() string {
	returnType, ok := p.expectType(NameTokenType)
	if !ok {
		fmt.Printf("ERROR: parseType: expected type %s\n",
			returnType)
		return ""
	}
	return string(returnType.type_)
}

func (p *Parser) parseFunc() *Func {
	_, ok := p.expectType(FuncTokenType)
	if !ok {
		//fmt.Printf("ERROR: parse func: expected type: %s\n", FuncTokenType)
		return nil
	}
	nameToken, ok := p.expectType(NameTokenType)
	if !ok {
		fmt.Printf("ERROR: parse func: expected type: %s\n", NameTokenType)
		return nil
	}

	params := p.parseParams()
	block := p.parseBlock()

	return &Func{
		name:   nameToken.value,
		body:   block,
		params: params,
	}
}

func (p *Parser) parseBlock() []Statement {
	if _, ok := p.expectType(ColonTokenType); !ok {
		return nil
	}

	block := make([]Statement, 0)
	for {
		token, ok := p.expectType(NameTokenType, EndTokenType, DefTokenType)
		if !ok {
			fmt.Printf("ERROR: parse block: expected types: %s, %s, %s, got: %s\n", NameTokenType, EndTokenType, DefTokenType, token.type_)
			break
		}
		end := false
		switch token.type_ {
		case EndTokenType:
			end = true
		case NameTokenType:
			params := p.parseArgs()
			block = append(block, FuncCall{
				name: token.value,
				args: params,
			})
		case DefTokenType:
			def, ok := p.parseDef()
			if !ok {
				panic("ERROR: parse block: def is empty")
			}
			block = append(block, def)
		}
		if end {
			break
		}
	}
	return block
}

func (p *Parser) parseDef() (Assignment, bool){
	name, ok := p.expectType(NameTokenType)
	if !ok {
		fmt.Printf("ERROR: parseDef: expected: %s, got: %s\n", NameTokenType,name.type_)
		return Assignment{}, false
	}
	eq, eqok := p.expectType(EqTokenType)
	if !eqok {
		fmt.Printf("ERROR: parseDef: expected: %s, got: %s\n", EqTokenType, eq.type_)
		return Assignment{}, false
	}
	valToken, ok := p.expectType(NumTokenType, StringTokenType)
	var ex Expression
	if ok {
		switch valToken.type_ {
		case NumTokenType:
			str, err := strconv.ParseInt(valToken.value, 10, 64)
			if err != nil {
				fmt.Printf("ERROR: parse def: cant conver int to str\n")
				return Assignment{}, false
			}
			ex = NumberLiteral{
				value: int(str),
			}
		case StringTokenType:
			ex = StringLiteral{
				value: valToken.value,
			}
		}
	}

	return Assignment{
		VarName: name.value,
		Value:   ex,
	}, true
}

func (p *Parser) parseParams() []string {
	_, ok := p.expectType(OparenTokenType)
	if !ok {
		return []string{}
	}

	params := []string{}

	for {
		token, ok := p.expectType(
			CparenTokenType,
			NameTokenType,
		)

		if !ok {
			return params
		}

		switch token.type_ {
		case CparenTokenType:
			return params

		case NameTokenType:
			params = append(params, token.value)
		}
	}
}

func (p *Parser) parseArgs() []Expression {
	_, ok := p.expectType(OparenTokenType)
	params := make([]Expression, 0)
	if !ok {
		return params
	}
	for {
		end := false
		if token, ok := p.expectType(CparenTokenType, NameTokenType, StringTokenType, CommaTokenType); ok {
			switch token.type_ {
			case CparenTokenType:
				end = true
			case StringTokenType:
				params = append(params, StringLiteral{
					value: token.value,
				})
			case NameTokenType:
				params = append(params, Variable{
					name: token.value,
				})
			case CommaTokenType:
				continue
			}
		}
		if end {
			break
		}
	}
	return params
}
