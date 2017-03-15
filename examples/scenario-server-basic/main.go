package main

import (
	"context"

	_ "github.com/moul/mainager/module/backend/raven"
	_ "github.com/moul/mainager/module/backend/smtp"
	"github.com/moul/mainager/scenario/server"
)

func main() {
	server_scenario.Run(context.Background())
}
