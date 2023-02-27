//go:build ignore
// +build ignore

package main

var a, b, c int

var unknown bool // defeat dead-code elimination

func func1() {
	var h int // @line f1h
	f := func(x *int) *int {
		if unknown {
			return &b
		}
		return x
	}

	// FV(g) = {f, h}
	g := func(x *int) *int {
		if unknown {
			return &h
		}
		return f(x)
	}

	print(g(&a)) // @pointsto command-line-arguments.a | command-line-arguments.b | h@f1h:6
	print(f(&a)) // @pointsto command-line-arguments.a | command-line-arguments.b
	print(&a)    // @pointsto command-line-arguments.a
}

var f func(*int) *int

func func2() {
	if unknown {
		f = func(x *int) *int {
			return x
		}
	} else {
		f = func(x *int) *int {
			return &c
		}
	}

	print(f(&a)) // @pointsto command-line-arguments.a | command-line-arguments.b
	print(f(&b))
}

func foo(bar func(*int) *int) *int {
	return bar(&b)
}

func func3() {
	f := func(x *int) *int {
		return x
	}

	print(foo(f))
	print(f(&a))
}

type S struct{ x *int }

func (s *S) foo(y *int) *int {
	return s.x
}

func func4() {
	var g func(*int) *int
	var s1 = &S{&a}
	//var s2 = &S{&b}
	g = (*S).foo(s1)
	print(g(&b))
	print(g(&c))

}

func main() {
	//func1()
	//func2()
	//func3()
	func4()
}
