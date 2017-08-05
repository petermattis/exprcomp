package exprcomp

import "unsafe"

type instruction byte

const (
	DATUM instruction = iota
	PLUS
	MINUS
)

type prog struct {
	ins  []instruction
	data []unsafe.Pointer
}

type evalContext struct {
	stack [16]unsafe.Pointer
}

func (p *prog) eval(ctx *evalContext) datum {
	stack := ctx.stack[:len(ctx.stack)]
	var sptr, dptr int
	for _, i := range p.ins {
		switch i {
		case DATUM:
		case PLUS:
			sptr -= 2
			*((*dint)(p.data[dptr])) = *((*dint)(stack[sptr+1])) + *((*dint)(stack[sptr]))
		case MINUS:
			sptr -= 2
			*((*dint)(p.data[dptr])) = *((*dint)(stack[sptr+1])) - *((*dint)(stack[sptr]))
		default:
			panic("not reached")
		}
		stack[sptr] = p.data[dptr]
		sptr++
		dptr++
	}
	return (*dint)(stack[sptr-1])
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
	p.data = append(p.data, unsafe.Pointer(newInt(0)))
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
	p.data = append(p.data, unsafe.Pointer(newInt(0)))
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
	p.data = append(p.data, unsafe.Pointer(n))
}
