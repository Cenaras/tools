//go:build ignore
// +build ignore

package main

var a, b int

type S struct{}

func (s *S) foo() *int {
	return &a
}

func staticInvoke1() {
	var s *S
	var x *int
	x = s.foo()
	print(x)
}

func main() {
	staticInvoke1()
}
