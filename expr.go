package exprcomp

import (
	"time"
)

type instruction byte

const (
	DATUM instruction = iota
	PLUS
	MINUS
)

type slot struct {
	i dint
	f float64
	s string
	t time.Time
}

type prog struct {
	ins  []instruction
	data []slot
}

type evalContext struct {
	stack [16]*slot
}

func (p *prog) eval(ctx *evalContext) datum {
	stack := ctx.stack[:len(ctx.stack)]
	var sptr, dptr int
	for _, i := range p.ins {
		switch i {
		case DATUM:
		case PLUS:
			sptr -= 2
			p.data[dptr].i = stack[sptr+1].i + stack[sptr].i
		case MINUS:
			sptr -= 2
			p.data[dptr].i = stack[sptr+1].i - stack[sptr].i
		default:
			panic("not reached")
		}
		stack[sptr] = &p.data[dptr]
		sptr++
		dptr++
	}
	return &stack[sptr-1].i
}

type expr interface {
	eval() datum
	compile(p *prog)
}

type datum interface {
	expr
}

type plus struct {
	left, right expr
}

func (n *plus) eval() datum {
	return newInt(int64(*n.left.eval().(*dint) + *n.right.eval().(*dint)))
}

func (n *plus) compile(p *prog) {
	n.right.compile(p)
	n.left.compile(p)
	p.ins = append(p.ins, PLUS)
	p.data = append(p.data, slot{})
}

type minus struct {
	left, right expr
}

func (n *minus) eval() datum {
	return newInt(int64(*n.left.eval().(*dint) - *n.right.eval().(*dint)))
}

func (n *minus) compile(p *prog) {
	n.right.compile(p)
	n.left.compile(p)
	p.ins = append(p.ins, MINUS)
	p.data = append(p.data, slot{})
}

type dint int64

func newInt(v int64) *dint {
	return (*dint)(&v)
}

func (n *dint) eval() datum {
	return n
}

func (n *dint) compile(p *prog) {
	p.ins = append(p.ins, DATUM)
	p.data = append(p.data, slot{i: *n})
}
