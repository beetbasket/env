package env

import (
	"fmt"
	"maps"
	"os"
	"reflect"
	"slices"
	"strings"
)

type parseFn func(_ regMap, _, _ string, _ ...string) (reflect.Value, bool)

func get[P KindParser[T], T any](registered regMap, key, help string, def ...string) (value reflect.Value, ok bool) {
	register(registered, key, help, def...)
	p := *new(P)
	if v, ok := os.LookupEnv(key); ok {
		vv, err := p.Parse(v)
		if err != nil {
			panic(fmt.Errorf("failed to parse env key %q, %v", key, err))
		}
		return reflect.ValueOf(vv), true
	} else if len(def) > 0 {
		vv, err := p.Parse(def[0])
		if err != nil {
			panic(fmt.Errorf("failed to parse default env key %q, %v", key, err))
		}
		return reflect.ValueOf(vv), true
	}
	return value, false
}

type regMap map[string]*info
type info struct {
	Type    string
	Help    string
	Default *string
}

func register[T any](registered regMap, key, help string, def ...T) {
	inf := info{
		Type: fmt.Sprintf("%T", *new(T)),
		Help: help,
	}

	if len(def) > 0 {
		d := fmt.Sprint(def)
		inf.Default = &d
	}

	registered[key] = &inf
}

func (rm regMap) help() string {
	var sb strings.Builder
	sb.WriteString("Environment:\n")
	for key := range slices.Values(slices.Sorted(maps.Keys(rm))) {
		inf := rm[key]
		_, _ = fmt.Fprintf(&sb, "  %s: %s %s", key, inf.Type, inf.Help)
		if inf.Default != nil {
			_, _ = fmt.Fprintf(&sb, " (default: %s)", *inf.Default)
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}
