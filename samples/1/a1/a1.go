package a1

type IA interface {
	Plus(int, int) int
	Minus(a, b int) (c int)
}

type B struct{}

func (b *B) Plus(a1, a2 int) int {
	return a1 + a2
}

func (b *B) Minus(a, d int) (c int) {
	c = a - d
	return
}

func NewRefIA() IA {
	return &B{}
}

type C struct{}

func (b C) Plus(a1, a2 int) int {
	return a1 + a2
}

func (b C) Minus(a, d int) (c int) {
	c = a - d
	return
}

func NewNoRefIA() IA {
	return C{}
}
