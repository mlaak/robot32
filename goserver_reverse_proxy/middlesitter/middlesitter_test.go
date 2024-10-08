package middlesitter

import (
	"strconv"
	"testing"
)

func TestObservableReadCloser(t *testing.T) {
	orc := NewObservableReadCloser()

	total := 0

	fun1 := func() {
		total = total + 1
	}
	fun2 := func() {
		total = total + 2
	}
	fun5 := func() {
		total = total + 5
	}

	orc.AddOnCloseFunc(0, fun1)
	orc.AddOnCloseFunc(0, fun2)
	orc.AddOnCloseFunc(0, fun5)

	orc.CallAllOnCloseFuncs(0)
	orc.CallAllOnCloseFuncs(0) //second time should do nothing

	if total != 8 {
		t.Fatalf("TestObservableReadCloser Fail 1 - expected 8 but got " + strconv.Itoa(total) + " Did a closer not get called or did get called twice?")
	}

	//TODO: write more tests
}

func TestMiddleSitter(t *testing.T) {
	//TODO: write tests
}
