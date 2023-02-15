package main

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	// Don't assert @pointsto(t) since its label contains a fragile line number.
	run_engine("invoke_sens.go")
	t.Fail()
}

type S struct {
	i int
}

func TestMap(t *testing.T) {
	m := make(map[S]bool)
	p := S{1}
	m[p] = true
	q := S{1}
	fmt.Print(m[q])
	t.Fail()
}
