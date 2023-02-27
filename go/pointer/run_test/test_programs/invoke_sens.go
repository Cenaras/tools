//go:build ignore
// +build ignore

package main

// Test of context-sensitive treatment of certain function calls,
// e.g. static calls to simple accessor methods.

var a, b int

var unknown bool

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
	t := &S{&b}
	print(bar(s))
	print(bar(t))
}

func bar(i I) *int {
	return i.foo()
}

func context2() {
	var s1 *S
	if unknown {
		s1 = &S{&a}
	} else {
		s1 = &S{&b}
	}
	s1.x = &b
	print(s1.foo())
	//print(s2.foo())
}

func main() {
	//context1()
	context2()

}
