package ping_service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/moul/mainager"
	_ "github.com/moul/mainager/module/server/http"
)

func init() {
	mainager.Register(mainager.Module{
		Name: "mainager.module.service.ping",
		Hooks: mainager.Hooks{
			"service-init": serviceInit,
		},
		Dependencies: []string{"mainager.module.server.http"},
	})
}

func serviceInit(ctx context.Context, params ...interface{}) (context.Context, error) {
	mux := ctx.Value(mainager.Key("mainager.module.server.http.mux")).(*http.ServeMux)

	mux.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	return ctx, nil
}
