package commands

type Registry struct {
	commands map[string]func(args []string) error
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]func(args []string) error),
	}
}

func (r *Registry) Register(name string, fn func(args []string) error) {
	r.commands[name] = fn
}

func (r *Registry) Execute(name string, args []string) error {
	fn, ok := r.commands[name]
	if !ok {
		return nil
	}
	return fn(args)
}
