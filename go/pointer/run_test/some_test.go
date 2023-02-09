package main

import "testing"

func Test(t *testing.T) {
	// Don't assert @pointsto(t) since its label contains a fragile line number.
	run_engine("callsite_dup.go")
}
