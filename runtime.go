package main

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
	return NotImplemented
}
