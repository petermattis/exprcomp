package exprcomp

import (
	"testing"
)

func makeExpr() expr {
	return &plus{
		left: &minus{
			left: &plus{
				left: &minus{
					left: &plus{
						left:  newInt(1),
						right: newInt(2),
					},
					right: &plus{
						left:  newInt(3),
						right: newInt(4),
					},
				},
				right: newInt(5),
			},
			right: newInt(6),
		},
		right: newInt(7),
	}
}

func BenchmarkEval(b *testing.B) {
	e := makeExpr()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if v := e.eval().(*dint); *v != 2 {
			b.Fatalf("unexpected value: %d", *v)
		}
	}
}

func BenchmarkCompile(b *testing.B) {
	e := makeExpr()
	var p prog
	e.compile(&p)
	stack := make([]datum, 0, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if v := p.eval(stack).(*dint); *v != 2 {
			b.Fatalf("unexpected value: %d", *v)
		}
	}
}
