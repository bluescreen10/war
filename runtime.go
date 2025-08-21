package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bluescreen10/war/text"
)

type FuncMap map[string]func(_, _ any)

type Runtime struct {
	globalFuncs FuncMap
}

type RuntimeOption func(*Runtime)

func NewRuntime(opts ...RuntimeOption) *Runtime {
	r := &Runtime{}
	for _, o := range opts {
		o(r)
	}
	return r
}

func WithFuncs(funcs FuncMap) RuntimeOption {
	return func(r *Runtime) {
		r.globalFuncs = funcs
	}
}

func (r *Runtime) ExecFile(path string) error {
	switch filepath.Ext(path) {
	case ".wat", ".wast":
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error opening file: %s", path)
		}

		p := text.NewParser(data)

		if err := p.Parse(); err != nil {
			return fmt.Errorf("parsing error: %v", err)
		}
		return nil
		//return t.Exec()
	default:
		return ErrNotImplemented
	}
}
