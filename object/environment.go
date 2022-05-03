package object

type Environment struct {
	store map[string]Object
	Outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.Outer = outer
	return env
}

func (e *Environment) GetAt(depth int, name string) (Object, bool) {
	obj, ok := e.ancestor(depth).store[name]
	return obj, ok
}

func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}

func (e *Environment) AssignAt(depth int, name string, value Object) {
	e.ancestor(depth).store[name] = value
}

func (e *Environment) ancestor(depth int) *Environment {
	env := e
	for i := 0; i < depth; i++ {
		env = env.Outer
	}

	return env
}

// func (e *Environment) Get(name string) (Object, bool) {
// 	obj, ok := e.store[name]
// 	if !ok && e.outer != nil {
// 		obj, ok = e.outer.Get(name)
// 	}
// 	return obj, ok
// }

// func (e *Environment) Assign(name string, value Object) bool {
// 	if _, ok := e.store[name]; ok {
// 		e.store[name] = value
// 		return true
// 	} else if e.outer != nil {
// 		return e.outer.Assign(name, value)
// 	} else {
// 		return false
// 	}
// }
