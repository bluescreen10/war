package main_test

import (
	"path/filepath"
	"testing"

	war "github.com/bluescreen10/war"
)

func TestSpec(t *testing.T) {
	matches, err := filepath.Glob(filepath.Join("testsuite", "*.wast"))
	if err != nil {
		t.Fatal("can't find test files")
	}

	for _, match := range matches {
		t.Run(match, func(t *testing.T) {
			runtime := NewTestRuntime(t)
			if err := runtime.ExecFile(match); err != nil {
				t.Errorf("runtime error: %v", err)
			}
		})
	}
}

func NewTestRuntime(t *testing.T) *war.Runtime {
	return war.NewRuntime(war.WithFuncs(war.FuncMap{
		"assert_return": func(got, expected any) {
			if expected != got {
				t.Errorf("assert_return: got(%v) expected(%v)", got, expected)
			}
		},
	}))
}
