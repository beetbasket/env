package env

import (
	"errors"
	"fmt"
	"reflect"
)

func Unmarshal(v any) (help func() string, _ error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer && rv.Kind() != reflect.Interface {
		return nil, errors.New("not a pointer to a value")
	}

	rv = rv.Elem()
	if !rv.IsValid() {
		return nil, errors.New("nil value")
	} else if rv.Kind() != reflect.Struct {
		return nil, errors.New("not a pointer to a struct")
	}

	rt := rv.Type()
	rm := regMap{}
	for i := 0; i < rv.NumField(); i++ {
		sv := rv.Field(i)
		st := rt.Field(i)
		if !st.IsExported() {
			return nil, errors.New("non exported struct field")
		}

		key, help, parser, def, err := unmarshalTags(st.Tag)
		if err != nil {
			return nil, err
		}

		var pFn parseFn
		switch parser {
		case "string":
			pFn = get[KindString, string]
		case "absFile":
			pFn = get[KindAbsFile, []string]
		case "args":
			pFn = get[KindArgs, []string]
		case "port":
			pFn = get[KindPort, uint16]
		default:
			return nil, fmt.Errorf("parser not found: %q", parser)
		}

		if sv.CanSet() {
			if v, ok := pFn(rm, key, help, def...); ok {
				sv.Set(v)
			}
		}
	}
	return rm.help, nil
}

type getFn func(key, help string, def ...string) reflect.Value

func unmarshalTags(tags reflect.StructTag) (key, help, parser string, def []string, err error) {
	var ok bool
	if key, ok = tags.Lookup("key"); !ok {
		err = errors.Join(err, errors.New("tag key not set"))
	}

	if help, ok = tags.Lookup("help"); !ok {
		err = errors.Join(err, errors.New("tag help not set"))
	}

	if parser, ok = tags.Lookup("parser"); !ok {
		err = errors.Join(err, errors.New("tag parser not set"))
	}

	if d, ok := tags.Lookup("default"); ok {
		def = []string{d}
	}
	return
}
