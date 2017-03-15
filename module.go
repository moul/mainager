package mainager

type Module struct {
	Name         string
	Hooks        Hooks
	Dependencies []string
}
