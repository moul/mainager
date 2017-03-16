package mainager

import "github.com/moul/mainager"

type Module struct {
	Name         string
	Hooks        Hooks
	Dependencies []string
	Context      mainager.Context
}

type Modules map[string]Module
