package mainager

import "context"

type Mainager struct {
	modules Modules
	context Context
}

func New() *Mainager {
	return &Mainager{
		modules: make(Modules, 0),
		context: make(Context, 0),
	}
}

func (m *Mainager) Register(module Module) {
	m.modules[module.Name] = module
}

func (m *Mainager) InvokeAll(ctx context.Context, hookName string, params ...interface{}) (context.Context, error) {
	for _, module := range m.modules {
		hook, found := module.Hooks[hookName]
		if !found {
			continue
		}

		var err error
		if ctx, err = hook(ctx, params...); err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}
