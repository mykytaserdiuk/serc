package core

import (
	"fmt"
	"strconv"

	"github.com/mykytaserdiuk/serc/ast"
)

type Parser struct {
	l       *Lexer
	current *ast.Token
	nextTok *ast.Token
}

func (p *Parser) next() *ast.Token {
	token := p.l.NextToken()
	//fmt.Println(token.Value)
	return token
}

// watch current token
func (p *Parser) peek() *ast.Token {
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

func (p *Parser) peekNext() *ast.Token {
	p.peek()

	if p.nextTok == nil {
		p.nextTok = p.next()
	}

	return p.nextTok
}

// take current token and go
func (p *Parser) advance() *ast.Token {
	t := p.peek()
	p.current = nil

	return t
}

// check token, but dont take him
func (p *Parser) match(types ...ast.TokenType) bool {
	token := p.peek()

	for _, t := range types {
		if token.Type_ == t {
			return true
		}
	}

	return false
}

// check token, and take him
func (p *Parser) expect(types ...ast.TokenType) (*ast.Token, bool) {
	token := p.peek()

	for _, t := range types {
		if token.Type_ == t {
			return p.advance(), true
		}
	}

	return token, false
}

func (p *Parser) parseProgram() ast.Program {
	program := ast.Program{
		Structs:    make(map[string]*ast.Structure),
		Functions:  make(map[string]*ast.Func),
		Imports:    make(map[string]ast.Import),
		BuildinFns: make(map[string]ast.BuiltinFunc),
	}
	token := p.peek()
	for token.Type_ != ast.EOFTokenType {
		//fmt.Println(token.Type_)
		switch token.Type_ {
		case ast.FuncTokenType:
			fn := p.parseFunc()
			program.Functions[fn.Name] = fn
		case ast.StructTokenType:
			str := p.parseStruct()
			program.Structs[str.Name] = str
		case ast.UseTokenType:
			use := p.parseUse()
			program.Imports[use.Name] = use
		default:
			p.advance()
		}
		token = p.peek()
	}
	return program
}

func (p *Parser) parseFunc() *ast.Func {
	_, ok := p.expect(ast.FuncTokenType)
	if !ok {
		//fmt.Printf("ERROR: parse func: expected type: %s\n", FuncTokenType)
		return nil
	}
	nameToken, ok := p.expect(ast.NameTokenType)
	if !ok {
		fmt.Printf("ERROR: parse func: expected type: %s\n", ast.NameTokenType)
		return nil
	}

	params, lastToken := p.parseParams()
	if lastToken.Type_ != ast.CparenTokenType {
		panic("ERROR: Expected close paren after params, but  got: " + lastToken.Type_)
	}
	if colon, ok := p.expect(ast.ColonTokenType); !ok {
		fmt.Printf("ERROR: parse func: expected type: %s, got %s\n", ast.ColonTokenType, colon.Type_)
		return nil
	}
	block, closeToken := p.parseBlock()
	if closeToken.Type_ != ast.EndTokenType {
		panic("ERROR: expected end token on end of function")
	}
	return &ast.Func{
		Name:   nameToken.Value,
		Body:   block.Statements,
		Params: params,
	}
}

func (p *Parser) parseBlock() (ast.Block, *ast.Token) {
	block := make([]ast.Statement, 0)
	var token *ast.Token
	for {
		token = p.peek()
		end := false
		switch token.Type_ {
		case ast.EndTokenType, ast.EOFTokenType, ast.ElseTokenType:
			return ast.Block{
				Statements: block,
			}, p.advance()
		case ast.NameTokenType:
			target := p.parseExpression() // name
			if p.match(ast.EqTokenType) {
				p.advance()
				value := p.parseExpression()

				block = append(block, ast.Assign{
					Target: target,
					Value:  value,
				})
			} else {
				switch e := target.(type) {
				case ast.Call:
					block = append(block, e)
				default:
					panic("unexpected variable statement: " + token.Value)
				}
			}
		case ast.WhileTokenType:
			loop := p.parseLoop()
			block = append(block, loop)
		case ast.ReturnTokenType:
			retToken := p.advance()
			ret := p.parseReturn(retToken)
			block = append(block, ret)
		case ast.DefTokenType:
			p.advance()
			def, ok := p.parseDef()
			if !ok {
				panic("ERROR: parse block: def is empty")
			}
			block = append(block, def)
		case ast.IfTokenType:
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
	return ast.Block{
		Statements: block,
	}, token
}

func (p *Parser) parseDef() (ast.NewAssign, bool) {
	name, ok := p.expect(ast.NameTokenType)
	if !ok {
		fmt.Printf("ERROR: parseDef: expected: %s, got: %s\n", ast.NameTokenType, name.Type_)
		return ast.NewAssign{}, false
	}
	eq, eqok := p.expect(ast.EqTokenType)
	if !eqok {
		fmt.Printf("ERROR: parseDef: expected: %s, got: %s\n", ast.EqTokenType, eq.Type_)
		return ast.NewAssign{}, false
	}

	exp := p.parseExpression()

	return ast.NewAssign{
		VarName: name.Value,
		Value:   exp,
	}, true
}

func (p *Parser) parseStruct() *ast.Structure {
	p.advance() // chop 'struct'

	name, ok := p.expect(ast.NameTokenType)
	if !ok {
		panic("expected name literal after 'struct'")
	}
	fields := p.parseStructFields()
	if len(fields) == 0 {
		panic("struct must have a fields")
	}
	return &ast.Structure{
		Name:   name.Value,
		Fields: fields,
	}
}

func (p *Parser) parseStructFields() []string {
	fields := []string{}
	_, ok := p.expect(ast.OparenTokenType)
	if !ok {

	}
	end := false
	for p.match(ast.CparenTokenType, ast.NameTokenType) {
		token := p.advance()
		switch token.Type_ {
		case ast.CparenTokenType:
			end = true
		case ast.NameTokenType:
			fields = append(fields, token.Value)
		}
		if end {
			break
		}
	}
	return fields
}

func (p *Parser) parseUse() ast.Import {
	p.advance() // use
	name := p.parsePrimary()
	var useName string
	var alias string
	switch t := name.(type) {
	case ast.Variable:
		useName = t.Name
	default:
		panic(fmt.Sprintf("expected name of use got: %T", t))
	}

	ok := p.match(ast.StringTokenType)
	if !ok {
		alias = useName
	} else {
		aliasExp := p.parsePrimary()
		if aliasVal, ok := aliasExp.(ast.StringLiteral); !ok {
			panic(fmt.Sprintf("as alias cab be only string literal, got: %T", aliasExp))
		} else {
			alias = aliasVal.Value
		}
	}

	return ast.Import{
		Alias: alias,
		Name:  useName,
	}
}

func (p *Parser) parseReturn(returnToken *ast.Token) ast.Return {
	if p.peek().Line > returnToken.Line {
		return ast.Return{
			Value: nil,
		}
	} else {
		return ast.Return{
			Value: p.parseExpression(),
		}
	}
}

func (p *Parser) parseLoop() ast.Loop {
	p.advance() // while
	conditions, ok := p.parseConditions()
	if !ok {
		panic("ERROR: parseLoop: cant parse contidions")
	}
	colon, ok := p.expect(ast.ColonTokenType)
	if !ok {
		panic("ERROR: parseLoop: expected ':' after conditions, got: " + colon.Type_)
	}
	block, _ := p.parseBlock()
	return ast.Loop{
		Body:       block,
		Conditions: conditions,
	}
}

func (p *Parser) parseIf() (ast.If, bool) {
	conditions, ok := p.parseConditions()
	if !ok {
		panic("ERROR: parseIf: cant parse conditions")
	}
	then, ok := p.expect(ast.ThenTokenType)
	if !ok {
		panic("ERROR: parseIf: 'then' not found, got: " + then.Type_)
	}
	thenBlock, lastChoped := p.parseBlock()
	if lastChoped == nil {
		fmt.Printf("ERROR: parseIf: expected some choped: %+v, but is nil\n", []ast.TokenType{ast.ElseTokenType, ast.EndTokenType})
		return ast.If{}, false
	}
	elseBlock := ast.Block{}
	if lastChoped.Type_ == ast.ElseTokenType {
		elseBlock, lastChoped = p.parseBlock()
		if lastChoped.Type_ != ast.EndTokenType {
			panic("expected end after else")
		}
	} else if lastChoped.Type_ != ast.EndTokenType {
		panic("expected end")
	}

	return ast.If{
		Then:       thenBlock,
		Else:       elseBlock,
		Conditions: conditions,
	}, true
}

func (p *Parser) parseConditions() (ast.Expression, bool) {
	if _, ok := p.expect(ast.OparenTokenType); !ok {
		return nil, false
	}

	expr := p.parseExpression()
	if expr == nil {
		return nil, false
	}

	if _, ok := p.expect(ast.CparenTokenType); !ok {
		return nil, false
	}

	return expr, true
}

func (p *Parser) parseParams() ([]string, *ast.Token) {
	token, ok := p.expect(ast.OparenTokenType)
	if !ok {
		return []string{}, token
	}

	params := []string{}

	for {
		token, ok := p.expect(
			ast.CparenTokenType,
			ast.NameTokenType,
		)

		if !ok {
			return params, token
		}

		switch token.Type_ {
		case ast.CparenTokenType:
			return params, token

		case ast.NameTokenType:
			params = append(params, token.Value)
		}
	}
}

func (p *Parser) parseArgs() []ast.Argument {
	var args []ast.Argument
	if !p.match(ast.OparenTokenType) {
		panic("expected '(' before args")
	}
	p.advance() // (

	for p.peek().Type_ != ast.CparenTokenType {
		current := p.peek()
		next := p.peekNext()
		if current.Type_ == ast.NameTokenType && next.Type_ == ast.ColonTokenType {
			name := p.advance() // name
			p.advance()         // :
			argExpr := p.parseExpression()
			args = append(args, ast.Argument{
				Name:  name.Value,
				Value: argExpr,
			})
		} else {
			exp := p.parseExpression()
			args = append(args, ast.Argument{
				Value: exp,
			})
		}

		if ok := p.match(ast.CommaTokenType); ok {
			p.advance()
			continue
		}
	}

	p.advance()

	return args
}

func (p *Parser) parsePrimary() ast.Expression {
	token, ok := p.expect(
		ast.NumTokenType,
		ast.StringTokenType,
		ast.NameTokenType,
	)
	if !ok {
		panic("parsePrimary: Expected: primary types, got: " + token.Type_)
	}
	var expr ast.Expression
	switch token.Type_ {
	case ast.StringTokenType:
		expr = ast.StringLiteral{
			Value: token.Value,
		}
	case ast.NumTokenType:
		str, err := strconv.ParseInt(token.Value, 10, 64)
		if err != nil {
			panic("ERROR: parseExpression: cant convert int")
		}

		expr = ast.NumberLiteral{
			Value: int(str),
		}
	case ast.NameTokenType:
		expr = ast.Variable{
			Name: token.Value,
		}

		for p.match(ast.DotTokenType) {
			p.advance() // .

			name, ok := p.expect(ast.NameTokenType)
			if !ok {
				panic("expected name after '.'")
			}

			expr = ast.FieldAccess{
				Value: expr,
				Name:  name.Value,
			}
		}

		if p.peek().Type_ == ast.OparenTokenType {
			args := p.parseArgs()
			expr = ast.Call{
				Target: expr,
				Args:   args,
			}
		}
	}

	return expr
}

func (p *Parser) parseAddition() ast.Expression {
	left := p.parsePrimary()
	for p.match(
		ast.PlusTokenType,
		ast.MinusTokenType,
	) {
		op := p.advance()

		right := p.parsePrimary()

		left = ast.Binary{
			Left:  left,
			Op:    op,
			Right: right,
		}
	}
	return left
}

func (p *Parser) parseComparation() ast.Expression {
	left := p.parseAddition()

	for p.match(
		ast.MoreTokenType,
		ast.LessTokenType,
		ast.EqEqTokenType,
		ast.EqMoreTokenType,
		ast.EqLessTokenType,
		ast.NotEqTokenType,
	) {
		op := p.advance()
		right := p.parseAddition()

		left = ast.Binary{
			Left:  left,
			Op:    op,
			Right: right,
		}
	}

	return left
}
func (p *Parser) parseExpression() ast.Expression {
	exp := p.parseComparation()
	return exp
}
