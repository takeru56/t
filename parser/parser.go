package parser

import (
	"strconv"

	"github.com/takeru56/t/token"
)

// Parser has the information of curToken and peekToken
type Parser struct {
	tokenizer *token.Tokenizer
	curToken  token.Token
	peekToken token.Token
}

// New initialize a Parser and returns its pointer
func New(t *token.Tokenizer) *Parser {
	p := &Parser{tokenizer: t}
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken advances forward curToken in the Parser
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.tokenizer.Next()
}

func (p *Parser) consume(s string) bool {
	if p.curToken.Literal == s {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) stmt() Node {
	return p.expr()
}

func (p *Parser) expr() Node {
	node := p.add()
	return node
}

func (p *Parser) add() Node {
	node := p.mul()
	tok := p.curToken
	for {
		if p.consume("+") {
			node = InfixExpr{tok, opAdd, node, p.mul()}
		} else if p.consume("-") {
			node = InfixExpr{tok, opSub, node, p.mul()}
		} else {
			return node
		}
	}
}

func (p *Parser) mul() Node {
	node := p.primary()
	tok := p.curToken
	for {
		if p.consume("*") {
			node = InfixExpr{tok, opMul, node, p.primary()}
		} else if p.consume("/") {
			node = InfixExpr{tok, opDiv, node, p.primary()}
		} else {
			return node
		}
	}
}

func (p *Parser) primary() Node {
	return p.newIntegerLiteral()
}

func (p *Parser) newIntegerLiteral() Node {
	val, _ := strconv.Atoi(p.curToken.Literal)
	node := IntegerLiteral{p.curToken, val}
	p.nextToken()
	return node
}
