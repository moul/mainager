package mainager

var instance *Mainager

func init() {
	instance = New()
}

func Instance() *Mainager {
	return instance
}

func Register(module Module) {
	instance.Register(module)
}
