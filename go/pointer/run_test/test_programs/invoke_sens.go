//go:build ignore
// +build ignore

package main

// Test of context-sensitive treatment of certain function calls,
// e.g. static calls to simple accessor methods.

var a, b int

type I interface {
	foo() *int
}

type T struct{ x *int }

func (t *T) foo() *int {
	return t.x
}

type S struct{ x *int }

func (s *S) foo() *int {
	return s.x
}

func context1() {
	s := &S{&a}
	t := &T{&b}
	print(bar(s))
	print(bar(t))
}

func bar(i I) *int {
	return i.foo()
}

func main() {
	context1()
}
