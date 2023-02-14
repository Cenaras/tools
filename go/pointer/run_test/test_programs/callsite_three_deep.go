//go:build ignore
// +build ignore

package main

func main() {
	var i, j, k int

	foo(&i)
	foo(&j)
	foo(&k)
}

func foo(x *int) *int {
	return bar4(x)
}

func bar4(x *int) *int {
	return baz(x)
}

func baz(x *int) *int {
	return x
}
