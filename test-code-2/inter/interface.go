package pkginter

type INTER interface {
	Plus(a, b int) (c int)
	Minus(int, int) int
	Repeat(a int, b int) (c, d int)
}
