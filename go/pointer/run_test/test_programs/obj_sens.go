//go:build ignore
// +build ignore

package main

// Test of context-sensitive treatment of certain function calls,
// e.g. static calls to simple accessor methods.

var a, b int

var unknown bool

type I interface {
	foo(x *int) *int
	bar()
}

type T struct{ x *int }

func (t *T) foo(x *int) *int {
	return x
}

type S struct {
	x *int
	y *int
}

func (s *S) foo(x *int) *int {
	return x
}

func (s *S) bar() {
	s.y = s.x
}

func context1() {
	var s I = &S{}
	print(s.foo(&a))
	print(s.foo(&b))
}

func context2() {
	var s I
	s1 := &S{x: &a}
	s2 := &S{x: &b}
	if unknown {
		s = s1
	} else {
		s = s2
	}
	s.bar()
	print(s1.y)
	print(s2.y)
}

func main() {
	//context1()
	context2()
}
