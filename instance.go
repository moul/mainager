package mainager

import "context"

var instance *Mainager

func init() {
	instance = New(context.Background())
}

func Instance() *Mainager {
	return instance
}
