package internal

import "testing"

func TestRetype(t *testing.T) {
	a := &ArgType{}
	t.Log(a.retype("[]int32"))
}
