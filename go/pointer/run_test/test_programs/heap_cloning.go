//go:build ignore
// +build ignore

package main

var a, b int

type S struct {
	x *int
}

func bar(x *int) *S {
	return &S{x}
}

func heap1() {
	s := bar(&a)
	t := bar(&b)
	print(s.x)
	print(t.x)
}

func main() {
	heap1()
}
