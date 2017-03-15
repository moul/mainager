package mainager

import (
	"context"
	"fmt"
)

type Mainager struct {
	modules []Module
	ctx     context.Context
}

func New(ctx context.Context) *Mainager {
	return &Mainager{
		modules: make([]Module, 0),
		ctx:     ctx,
	}
}

func (m *Mainager) Register(module Module) {
	m.modules = append(m.modules, module)
}

func (m *Mainager) InvokeAll(hook string, params ...interface{}) error {
	return fmt.Errorf("not implemented")
}

func (m *Mainager) Invoke(module, hook string, params ...interface{}) error {
	return fmt.Errorf("not implemented")
}
