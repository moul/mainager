package mainager

type Module struct {
	Name         string
	Hooks        Hooks
	Dependencies []string
}

type Modules map[string]Module
