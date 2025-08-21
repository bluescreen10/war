package text

import (
	"errors"
	"fmt"
)

var ErrInvalidInput = errors.New("invalid input")

type Op int

const (
	OpUnkown Op = iota
	OpStart
	OpConst
	OpLocalGet
	OpLocalSet
	OpCall
	OpI32Add
)

var idCounter int

func newID() int {
	idCounter++
	return idCounter
}

type Node struct {
	ID   int
	Op   Op
	Args []*Node // inputs
	Meta string  // e.g. immediate value, func name
}

func NewNode(op Op, meta string, args ...*Node) *Node {
	return &Node{ID: newID(), Op: op, Meta: meta, Args: args}
}

type Parser struct {
	lex  *lexer
	root *Node
}

func NewParser(input []byte) *Parser {
	return &Parser{
		lex: NewLexer(input),
	}
}

func (p *Parser) Parse() error {
	p.root = NewNode(OpStart, "", nil)
	for {
		t := p.lex.nextToken()

		//fmt.Printf("token: %s\n", t)
		if t.kind == tokenEOF {
			break
		}

		if t.kind == tokenError {
			return fmt.Errorf("lexing error: %q", t.val)
		}
	}
	return nil
}
