package main

import (
	"context"

	_ "github.com/moul/mainager/module/service/ping-http"
	"github.com/moul/mainager/scenario/server"
)

func main() {
	server_scenario.Run(context.Background())
}
