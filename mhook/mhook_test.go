package mhook

import (
	"github.com/sudachen/misc/out"
	"testing"
)

var scallBind, _ = NewHook("ws2_32!bind#0")

func trace(hook Hook, a []uintptr) Result {
	out.StdErr.Printf("hook: %s", hook.Name())
	return Result{Continue: true}
}

func Test_Network(t *testing.T) {
	if err := scallBind.Set(trace); err != nil {
		t.Fatal(err)
	}
}
