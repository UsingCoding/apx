package sandbox

var (
	R = Registry{
		sandboxes: map[string]Sandbox{},
	}
)

type Registry struct {
	sandboxes map[string]Sandbox
}

func (r Registry) Register(id string, s Sandbox) {
	_, ok := r.sandboxes[id]
	if ok {
		panic("sandbox already registered: " + id)
	}

	r.sandboxes[id] = s
}

func (r Registry) Lookup(id string) (Sandbox, bool) {
	s, ok := r.sandboxes[id]
	return s, ok
}
