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
	return bar(x)
}

func bar(x *int) *int {
	return x
}
