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
	l       *Lexer
	current *Token
	nextTok *Token
}

func (p *Parser) next() *Token {
	token := p.l.NextToken()
	// fmt.Println(token.value)
	return token
}

// watch current token
func (p *Parser) peek() *Token {
	if p.current == nil {
		if p.nextTok != nil {
			p.current = p.nextTok
			p.nextTok = nil
		} else {
			p.current = p.next()
		}
	}

	return p.current
}

func (p *Parser) peekNext() *Token {
	if p.nextTok == nil {
		p.nextTok = p.next()
	}

	return p.nextTok
}

// take current token and go
func (p *Parser) advance() *Token {
	token := p.peek()
	p.current = nil
	return token
}

// check token, but dont take him
func (p *Parser) match(types ...TokenType) bool {
	token := p.peek()

	for _, t := range types {
		if token.type_ == t {
			return true
		}
	}

	return false
}

// check token, and take him
func (p *Parser) expect(types ...TokenType) (*Token, bool) {
	token := p.peek()

	for _, t := range types {
		if token.type_ == t {
			p.current = nil
			return token, true
		}
	}

	return token, false
}

func (p *Parser) parseFunc() *Func {
	_, ok := p.expect(FuncTokenType)
	if !ok {
		//fmt.Printf("ERROR: parse func: expected type: %s\n", FuncTokenType)
		return nil
	}
	nameToken, ok := p.expect(NameTokenType)
	if !ok {
		fmt.Printf("ERROR: parse func: expected type: %s\n", NameTokenType)
		return nil
	}

	params, lastToken := p.parseParams()
	if lastToken.type_ != CparenTokenType {
		panic("ERROR: Expected close paren after params, but  got: " + lastToken.type_)
	}
	if colon, ok := p.expect(ColonTokenType); !ok {
		fmt.Printf("ERROR: parse func: expected type: %s, got %s\n", ColonTokenType, colon.type_)
		return nil
	}
	block, closeToken := p.parseBlock()
	if closeToken.type_ != EndTokenType {
		panic("ERROR: expected end token on end of function")
	}
	return &Func{
		name:   nameToken.value,
		body:   block.Statements,
		params: params,
	}
}

func (p *Parser) parseBlock() (Block, *Token) {
	block := make([]Statement, 0)
	var token *Token
	for {
		token = p.peek()
		end := false
		switch token.type_ {
		case EndTokenType, EOFTokenType, ElseTokenType:
			return Block{
				Statements: block,
			}, p.advance()
		case NameTokenType:
			nextToken := p.peekNext()
			if nextToken.type_ == EqTokenType {
				p.advance()
				p.advance()
				value := p.parseExpression()

				block = append(block, Assign{
					VarName: token.value,
					Value:   value,
				})
			} else if nextToken.type_ == OparenTokenType {
				p.advance()
				p.advance()
				args := p.parseArgs()
				block = append(block, FuncCall{
					name: token.value,
					args: args,
				})
			} else {
				panic("unexpected variable statement: " + token.value)
			}
		case ReturnTokenType:
			retToken := p.advance()
			ret := p.parseReturn(retToken)
			block = append(block, ret)
		case DefTokenType:
			p.advance()
			def, ok := p.parseDef()
			if !ok {
				panic("ERROR: parse block: def is empty")
			}
			block = append(block, def)
		case IfTokenType:
			p.advance()
			ifStmt, ok := p.parseIf()
			if ok {
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

func (p *Parser) parseDef() (NewAssign, bool) {
	name, ok := p.expect(NameTokenType)
	if !ok {
		fmt.Printf("ERROR: parseDef: expected: %s, got: %s\n", NameTokenType, name.type_)
		return NewAssign{}, false
	}
	eq, eqok := p.expect(EqTokenType)
	if !eqok {
		fmt.Printf("ERROR: parseDef: expected: %s, got: %s\n", EqTokenType, eq.type_)
		return NewAssign{}, false
	}
	var exp Expression
	if ok {
		exp = p.parseExpression()
	}

	return NewAssign{
		VarName: name.value,
		Value:   exp,
	}, true
}

func (p *Parser) parseReturn(returnToken *Token) Return {
	if p.peek().line > returnToken.line {
		return Return{
			Value: nil,
		}
	} else {
		return Return{
			Value: p.parseExpression(),
		}
	}
}

func (p *Parser) parseIf() (If, bool) {
	conditions, ok := p.parseIfConditions()
	if !ok {
		panic("ERROR: parseIf: cant parse conditions")
	}
	then, ok := p.expect(ThenTokenType)
	if !ok {
		panic("ERROR: parseIf: 'then' not found, got: " + then.type_)
	}
	thenBlock, lastChoped := p.parseBlock()
	if lastChoped == nil {
		fmt.Printf("ERROR: parseIf: expected some choped: %+v, but is nil\n", []TokenType{ElseTokenType, EndTokenType})
		return If{}, false
	}
	elseBlock := Block{}
	if lastChoped.type_ == ElseTokenType {
		elseBlock, lastChoped = p.parseBlock()
		if lastChoped.type_ != EndTokenType {
			panic("expected end after else")
		}
	} else if lastChoped.type_ != EndTokenType {
		panic("expected end")
	}

	return If{
		Then:       thenBlock,
		Else:       elseBlock,
		Conditions: conditions,
	}, true
}

func (p *Parser) parseIfConditions() (Expression, bool) {
	if _, ok := p.expect(OparenTokenType); !ok {
		return nil, false
	}

	expr := p.parseExpression()
	if expr == nil {
		return nil, false
	}

	if _, ok := p.expect(CparenTokenType); !ok {
		return nil, false
	}

	return expr, true
}

func (p *Parser) parseParams() ([]string, *Token) {
	token, ok := p.expect(OparenTokenType)
	if !ok {
		return []string{}, token
	}

	params := []string{}

	for {
		token, ok := p.expect(
			CparenTokenType,
			NameTokenType,
		)

		if !ok {
			return params, token
		}

		switch token.type_ {
		case CparenTokenType:
			return params, token

		case NameTokenType:
			params = append(params, token.value)
		}
	}
}

func (p *Parser) parseArgs() []Expression {
	var args []Expression
	if ok := p.match(CparenTokenType); ok {
		p.advance()
		return args
	}

	for {
		arg := p.parseExpression()
		args = append(args, arg)

		if ok := p.match(CommaTokenType); ok {
			p.advance()
			continue
		}

		if p.match(CparenTokenType) {
			p.advance()
			break
		}
	}

	return args
}

func (p *Parser) parsePrimary() Expression {
	token, ok := p.expect(
		NumTokenType,
		StringTokenType,
		NameTokenType,
	)
	if !ok {
		panic("parsePrimary: Expected: primary types, got: " + token.type_)
	}
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
		if ok := p.match(OparenTokenType); ok {
			p.advance()
			args := p.parseArgs()
			return FuncCall{
				name: token.value,
				args: args,
			}
		} else {
			return Variable{
				name: token.value,
			}
		}
	}

	panic("ERROR: unsupported expression type: " + token.type_)
}

func (p *Parser) parseAddition() Expression {
	left := p.parsePrimary()
	for p.match(
		PlusTokenType,
		MinusTokenType,
	) {
		op := p.advance()

		right := p.parsePrimary()
		left = Binary{
			Left:  left,
			Op:    op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) parseComparation() Expression {
	left := p.parseAddition()
	for p.match(
		MoreTokenType,
		LessTokenType,
		EqEqTokenType,
		EqMoreTokenType,
		EqLessTokenType,
	) {
		op := p.advance()

		right := p.parseAddition()
		left = Binary{
			Left:  left,
			Op:    op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) parseExpression() Expression {
	return p.parseComparation()
}
