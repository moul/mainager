package http_server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/moul/mainager"
	"github.com/urfave/cli"
)

func init() {
	mainager.Register(mainager.Module{
		Name: "github.com/moul/mainager/module/server/http",
		Hooks: mainager.Hooks{
			"cli-init":     cliInit,
			"cli-parse":    cliParse,
			"server-init":  serverInit,
			"server-start": serverStart,
		},
	})
}

func cliInit(ctx context.Context, params ...interface{}) (context.Context, error) {
	if len(params) != 1 {
		return ctx, fmt.Errorf("not enough arguments")
	}
	app := params[0].(*cli.App)
	app.Flags = append(app.Flags, []cli.Flag{
		cli.StringFlag{
			Name:   "http-bind",
			Usage:  "HTTP bind address. If empty, HTTP will be disabled",
			EnvVar: "HTTP_BIND",
			Value:  ":8000",
		},
	}...)
	return ctx, nil
}

func cliParse(ctx context.Context, params ...interface{}) (context.Context, error) {
	if len(params) != 1 {
		return ctx, fmt.Errorf("not enough arguments")
	}
	c := params[0].(*cli.Context)
	return context.WithValue(ctx, mainager.Key("github.com/moul/mainager/module/server/http.settings.bind-address"), c.String("http-bind")), nil
}

func serverInit(ctx context.Context, params ...interface{}) (context.Context, error) {
	mux := http.NewServeMux()
	return context.WithValue(ctx, mainager.Key("github.com/moul/mainager/module/server/http.mux"), mux), nil
}

func serverStart(ctx context.Context, params ...interface{}) (context.Context, error) {
	if len(params) != 1 {
		return ctx, fmt.Errorf("not enough arguments")
	}
	errc := params[0].(chan error)

	mux := ctx.Value(mainager.Key("github.com/moul/mainager/module/server/http.mux")).(*http.ServeMux)
	address := ctx.Value(mainager.Key("github.com/moul/mainager/module/server/http.settings.bind-address")).(string)
	if address == "" {
		return ctx, nil
	}

	log.Printf("github.com/moul/mainager/module/server/http: listen %q", address)
	go func() {
		errc <- http.ListenAndServe(address, handlers.LoggingHandler(os.Stderr, mux))
	}()

	return ctx, nil
}
