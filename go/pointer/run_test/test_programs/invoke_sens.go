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

func (s S) foo() *int {
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
	var s1 I
	if unknown {
		s1 = S{&a}
	} else {
		s1 = S{&b}
	}
	//s1.x = &b
	print(s1.foo())
	//print(s2.foo())
}

type J interface {
	f()
}

type C int

func (*C) f() {}

type D struct{ ptr *int }

func (D) f() {}

type E struct{}

func (*E) f() {}

func interface2() {
	var i J = (*C)(&a)
	var j J = D{&a}
	k := j
	if unknown {
		k = i
	}

	print(i) // @types *C
	print(j) // @types D
	print(k) // @types *C | D
	print(k) // @pointsto makeinterface:command-line-arguments.D | makeinterface:*command-line-arguments.C

	k.f()
	// @calls command-line-arguments.interface2 -> (*command-line-arguments.C).f
	// @calls command-line-arguments.interface2 -> (command-line-arguments.D).f

	print(i.(*C))    // @pointsto command-line-arguments.a
	print(j.(D).ptr) // @pointsto command-line-arguments.a
	print(k.(*C))    // @pointsto command-line-arguments.a

	switch x := k.(type) {
	case *C:
		print(x) // @pointsto command-line-arguments.a
	case D:
		print(x.ptr) // @pointsto command-line-arguments.a
	case *E:
		print(x) // @pointsto
	}
}

func main() {
	//context1()
	//context2()
	interface2()
}
