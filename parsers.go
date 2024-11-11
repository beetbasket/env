package env

import (
	"errors"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type KindParser[T any] interface {
	Parse(string) (T, error)
}

type KindString struct{}

func (KindString) Parse(s string) (string, error) { return s, nil }

type KindAbsFile struct{}

func (KindAbsFile) Parse(s string) ([]string, error) {
	s, err := filepath.Abs(s)
	if err != nil {
		return nil, err
	}
	return splitPath(s), nil
}

func splitPath(path string) []string {
	dir, last := filepath.Split(filepath.Clean(path))
	if last == "" || last == "." {
		if dir != "" {
			return []string{dir}
		}
		return []string{}
	}
	return append(splitPath(dir), last)
}

type KindArgs struct{}

func (KindArgs) Parse(s string) ([]string, error) {
	const pattern = `("[^"]*"|[^"\s]+)(\s+|$)`
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	ss := r.FindAllString(s, -1)
	out1 := make([]string, len(ss))
	for i, v := range ss {
		out1[i] = strings.TrimSpace(v)
	}
	return out1, nil
}

type KindPort struct{}

func (KindPort) Parse(s string) (uint16, error) {
	i, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return 0, err
	} else if i == 0 {
		return 0, errors.New("port cannot be 0")
	}
	return uint16(i), nil
}
