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
		Name: "github.com/moul/mainager/module/service/ping-http",
		Hooks: mainager.Hooks{
			"service-init": serviceInit,
		},
		Dependencies: []string{"github.com/moul/mainager/module/server/http"},
	})
}

func serviceInit(ctx context.Context, params ...interface{}) (context.Context, error) {
	mux := ctx.Value(mainager.Key("github.com/moul/mainager/module/server/http.mux")).(*http.ServeMux)

	mux.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	return ctx, nil
}
