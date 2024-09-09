package main

import (
	"path"
	"strings"
)

func pkgPrefix(s ...string) func(string) string {
	for i, v := range s {
		s[i] = path.Base(v)
	}

	return func(m string) string {
		return "[" + strings.Join(s, "/") + "] " + m
	}
}
