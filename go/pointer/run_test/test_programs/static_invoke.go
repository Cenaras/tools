//go:build ignore
// +build ignore

package main

var a, b int

type S struct{}

func (s *S) foo() *int {
	return &a
}

func staticInvoke1(s *S) {
	var x *int
	if s != nil {
		x = s.foo()
	} else {
		x = &b
	}
	print(x)
}

func main() {
	staticInvoke1(nil)
}
