package sandbox

var (
	R = Registry{
		sandboxes: map[string]Sandbox{},
	}
)

type Registry struct {
	sandboxes map[string]Sandbox
}

func (r Registry) Register(s Sandbox) {
	_, ok := r.sandboxes[s.ID]
	if ok {
		panic("sandbox already registered: " + s.ID)
	}

	r.sandboxes[s.ID] = s
}

func (r Registry) Lookup(id string) (Sandbox, bool) {
	s, ok := r.sandboxes[id]
	return s, ok
}
