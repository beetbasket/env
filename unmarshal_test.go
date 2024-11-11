package env

import (
	"testing"
)

func TestUnmarshal(t *testing.T) {
	var v struct {
		Goroot []string `key:"GOROOT" help:"go root" parser:"absFile"`
		Foo    string   `key:"FOO" help:"foo" default:"bar" parser:"string"`
	}

	help, err := Unmarshal(&v)
	if err != nil {
		t.Errorf("error: %v", err)
	} else if help() == "" {
		t.Errorf("no help")
	} else if v.Foo != "bar" {
		t.Errorf("foo not correct")
	} else if len(v.Goroot) == 0 {
		t.Errorf("goroot empty")
	}
	t.Logf("GOROOT: %s", v.Goroot)
	t.Logf("FOO: %s", v.Foo)
	t.Log(help())
}
