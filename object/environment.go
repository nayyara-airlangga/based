package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

func NewLocalEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, exists := e.store[name]
	if !exists && e.outer != nil {
		obj, exists = e.outer.store[name]
	}
	return obj, exists
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
