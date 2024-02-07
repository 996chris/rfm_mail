package main

import (
	"os"
)

func EnvGet(key string, defal string) string {
	s := os.Getenv(key)
	if s == "" {
		s = defal
	}
	return s
}
