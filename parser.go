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
	if colon, ok := p.expectType(ColonTokenType); !ok {
		fmt.Printf("ERROR: parse func: expected type: %s, got %s\n", ColonTokenType, colon.type_)
		return nil
	}
	block,_ := p.parseBlock()

	return &Func{
		name:   nameToken.value,
		body:   block.Statements,
		params: params,
	}
}

func (p *Parser) parseBlock() (Block, *Token) {
	block := make([]Statement, 0)
	var token *Token
	var ok bool
	for {
		token, ok = p.expectType(NameTokenType, EndTokenType, ElseTokenType, DefTokenType, IfTokenType)
		if !ok {
			fmt.Printf("ERROR: parse block: expected types: %+v,  got: %s\n", []TokenType{NameTokenType, EndTokenType, DefTokenType, ElseTokenType}, token.type_)
			break
		}
		end := false
		switch token.type_ {
		case EndTokenType, ElseTokenType:
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
		case IfTokenType:
			ifStmt, ok := p.parseIf()
			if ok{
				block = append(block, ifStmt)
			} else {
				panic("CANT PARSE IF")
			}
		}
		if end {
			break
		}
	}
	return Block{
		Statements: block,
	}, token
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
	var exp Expression
	if ok {
		exp = p.parseExpression(valToken)
	}

	return Assignment{
		VarName: name.value,
		Value:   exp,
	}, true
}

func (p *Parser) parseIf() (If, bool) {
	conditions, ok := p.parseIfConditions()
	if !ok {
		panic("ERROR: parseIf: cant parse conditions")
	}
	then, ok := p.expectType(ThenTokenType)
	if !ok {
		panic("ERROR: parseIf: 'then' not found, got: " + then.type_)
	}
	thenBlock, lastChoped := p.parseBlock()
	if lastChoped == nil{
		fmt.Printf("ERROR: parseIf: expected some choped: %+v, but is nil\n", []TokenType{ElseTokenType, EndTokenType})
		return If{}, false
	}
	elseBlock := Block{}
	switch lastChoped.type_{
		case ElseTokenType:
		elseBlock,_ = p.parseBlock()
		case EndTokenType:
		break;
	}

	return If{
		Then: thenBlock,
		Else: elseBlock,
		Conditions: conditions,
	}, true
}

func (p *Parser) parseIfConditions() (Binary, bool) {
	open, ok := p.expectType(OparenTokenType)
	if !ok {
		fmt.Printf("ERROR: parseIfCond: expected: %s, got: %s\n", OparenTokenType, open.type_)
		return Binary{}, false
	}

	left, ok := p.expectType(NumTokenType, StringTokenType, NameTokenType)
	if !ok {
		fmt.Printf("ERROR: parseIfCond: expected: %v, got: %s\n", []TokenType{NumTokenType, StringTokenType, NameTokenType}, open.type_)
		return Binary{}, false
	}

	leftEx := p.parseExpression(left)
	op, ok := p.parseLogicOperator("parseIfConditions")
	if !ok {
		return Binary{}, false
	}

	right, ok := p.expectType(NumTokenType, StringTokenType, NameTokenType)
	if !ok {
		fmt.Printf("ERROR: parseIfCond: expected: %v, got: %s\n", []TokenType{NumTokenType, StringTokenType, NameTokenType}, open.type_)
		return Binary{}, false
	}
	rightEx := p.parseExpression(right)
	close, ok := p.expectType(CparenTokenType)
	if !ok {
		fmt.Printf("ERROR: parseIf: expected: %s, got: %s\n", CparenTokenType, close.type_)
		return Binary{}, false
	}

	return Binary{
		Left: leftEx,
		Op: op,
		Right: rightEx,
	}, true
}

func (p *Parser) parseLogicOperator(fnName string) (*Token, bool) {
	op, ok := p.expectType(MoreTokenType, LessTokenType, EqMoreTokenType, EqLessTokenType)
	if !ok {
		fmt.Printf("ERROR: %s: expected: %v, got: %s\n", fnName, []TokenType{MoreTokenType, LessTokenType, EqMoreTokenType, EqLessTokenType}, op.type_)
		return op, false
	}

	return op, ok
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
			case StringTokenType, NameTokenType:
				ex :=p.parseExpression(token)
				if ex == nil{
					panic("ERROR: cant parse expression from: " + token.type_)
				}
				params = append(params, ex)
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

func (p *Parser) parseExpression(token *Token) Expression {
	switch token.type_ {
	case StringTokenType:
		return StringLiteral{
			value: token.value,
		}
	case NumTokenType:
		str, err := strconv.ParseInt(token.value, 10, 64)
			if err != nil {
				fmt.Printf("ERROR: parseExpression: cant conver int to str\n")
				return nil
			}
			return NumberLiteral{
				value: int(str),
			}
	case NameTokenType:
		return  Variable{
			name: token.value,
		}
	}

	fmt.Printf("ERROR: unsupported expression type: %s", token.type_)
	return nil
}
