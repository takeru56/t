package parser

import (
	"strconv"

	"github.com/takeru56/tcompiler/token"
)

// Parser has the information of curToken and peekToken
type Parser struct {
	tokenizer *token.Tokenizer
	curToken  token.Token
	peekToken token.Token
}

// type ParseErr struct {
// 	Err error
// 	L   token.Loc
// }

// // custom error
// var (
// 	ErrSyntax = errors.New("Syntax error: unexpected token")
// )

// func (pe *ParseErr) Error() string {
// 	switch pe.Err {
// 	case ErrSyntax:
// 		return fmt.Sprintf("%v", pe.Err)
// 	}
// 	return pe.Err.Error()
// }

// New initialize a Parser and returns its pointer
func New(t *token.Tokenizer) (*Parser, error) {
	p := &Parser{tokenizer: t}
	err := p.nextToken()
	if err != nil {
		return p, err
	}
	err = p.nextToken()
	if err != nil {
		return p, err
	}
	return p, nil
}

// nextToken advances forward curToken in the Parser
func (p *Parser) nextToken() error {
	p.curToken = p.peekToken
	t, err := p.tokenizer.Next()
	if err != nil {
		return err
	}
	p.peekToken = t
	return nil
}

func (p *Parser) consume(s string) (bool, error) {
	if p.curToken.Literal == s {
		err := p.nextToken()
		if err != nil {
			return true, err
		}
		return true, nil
	}
	return false, nil
}

func (p *Parser) Program() ([]Node, error) {
	program := []Node{}
	for p.curToken.Kind != token.EOF {
		n, err := p.stmt()
		if err != nil {
			return nil, err
		}
		program = append(program, n)
	}
	return program, nil
}

func (p *Parser) stmt() (Node, error) {
	f, err := p.consume("if")
	if err != nil {
		return IfStmt{}, err
	}
	if f {
		block := BlockStmt{Nodes: []Node{}}
		node, err := p.expr()
		if err != nil {
			return node, err
		}

		_, err = p.consume("then")
		if err != nil {
			return IfStmt{}, err
		}

		for {
			f, err = p.consume("end")
			if err != nil {
				return IfStmt{}, err
			}
			if f || p.curToken.Kind == token.EOF {
				break
			}

			n, err := p.stmt()
			if err != nil {
				return IfStmt{}, err
			}
			block.Nodes = append(block.Nodes, n)
		}
		return IfStmt{Condition: node, Block: block}, nil
	}

	f, err = p.consume("while")
	if err != nil {
		return WhileStmt{}, err
	}
	if f {
		block := BlockStmt{Nodes: []Node{}}
		node, err := p.expr()
		if err != nil {
			return BlockStmt{}, err
		}
		_, err = p.consume("do")
		if err != nil {
			return BlockStmt{}, err
		}

		for {
			f, err = p.consume("end")
			if err != nil {
				return WhileStmt{}, err
			}
			if f || p.curToken.Kind == token.EOF {
				break
			}

			n, err := p.stmt()
			if err != nil {
				return WhileStmt{}, err
			}
			block.Nodes = append(block.Nodes, n)
		}
		return WhileStmt{Condition: node, Block: block}, nil
	}
	node, err := p.assign()
	if err != nil {
		return node, err
	}
	return node, nil
}

func (p *Parser) assign() (Node, error) {
	node, err := p.expr()
	if err != nil {
		return node, err
	}
	switch node.(type) {
	case IdentExpr:
		f, err := p.consume("=")
		if err != nil {
			return AssignStmt{}, err
		}
		if f {
			n, err := p.expr()
			if err != nil {
				return AssignStmt{}, err
			}
			return AssignStmt{node.(IdentExpr), n}, nil
		}
	}
	return node, nil
}

func (p *Parser) expr() (Node, error) {
	node, err := p.eq()
	if err != nil {
		return node, err
	}
	return node, nil
}

func (p *Parser) eq() (Node, error) {
	node, err := p.compare()
	if err != nil {
		return node, err
	}
	tok := p.curToken
	for {
		if f, err := p.consume("=="); f {
			if err != nil {
				return node, err
			}
			n, err := p.compare()
			if err != nil {
				return n, err
			}
			node = InfixExpr{tok, EQ, node, n}
		} else if f, err := p.consume("!="); f {
			if err != nil {
				return node, err
			}
			n, err := p.compare()
			if err != nil {
				return n, err
			}
			node = InfixExpr{tok, NEQ, node, n}
		} else {
			return node, nil
		}
	}
}

func (p *Parser) compare() (Node, error) {
	node, err := p.add()
	if err != nil {
		return node, err
	}
	tok := p.curToken
	for {
		if f, err := p.consume("<"); f {
			if err != nil {
				return node, err
			}
			n, err := p.add()
			if err != nil {
				return node, err
			}
			node = InfixExpr{tok, Less, node, n}
		} else if f, err := p.consume(">"); f {
			if err != nil {
				return node, err
			}
			n, err := p.add()
			if err != nil {
				return node, err
			}
			node = InfixExpr{tok, Greater, node, n}
		} else {
			return node, nil
		}
	}
}

func (p *Parser) add() (Node, error) {
	node, err := p.mul()
	if err != nil {
		return node, err
	}
	tok := p.curToken
	for {
		if f, err := p.consume("+"); f {
			if err != nil {
				return node, err
			}
			n, err := p.mul()
			if err != nil {
				return node, err
			}
			node = InfixExpr{tok, Add, node, n}
		} else if f, err := p.consume("-"); f {
			if err != nil {
				return node, err
			}
			n, err := p.mul()
			if err != nil {
				return node, err
			}
			node = InfixExpr{tok, Sub, node, n}
		} else {
			return node, nil
		}
	}
}

func (p *Parser) mul() (Node, error) {
	node := p.prim()
	tok := p.curToken
	for {
		if f, err := p.consume("*"); f {
			if err != nil {
				return node, err
			}
			node = InfixExpr{tok, Mul, node, p.prim()}
		} else if f, err := p.consume("/"); f {
			if err != nil {
				return node, err
			}
			node = InfixExpr{tok, Div, node, p.prim()}
		} else {
			return node, nil
		}
	}
}

// prim ::= atom |
func (p *Parser) prim() Node {
	return p.atom()
}

// atom ::= IntegerLiteral | Identifier
func (p *Parser) atom() Node {
	switch p.curToken.Kind {
	case token.Num:
		return p.newIntegerLiteral()
	}
	return p.newIdentifier()
}

func (p *Parser) newIntegerLiteral() Node {
	val, _ := strconv.Atoi(p.curToken.Literal)
	node := IntegerLiteral{p.curToken, val}
	p.nextToken()
	return node
}

func (p *Parser) newIdentifier() Node {
	node := IdentExpr{variable, p.curToken.Literal}
	p.nextToken()
	return node
}