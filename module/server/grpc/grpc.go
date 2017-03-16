package grpc_server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/moul/mainager"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func init() {
	mainager.Register(mainager.Module{
		Name: "github.com/moul/mainager/module/server/grpc",
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
			Name:   "grpc-bind",
			Usage:  "gRPC bind address. If empty, gRPC will be disabled",
			EnvVar: "GRPC_BIND",
			Value:  ":9000",
		},
	}...)
	return ctx, nil
}

func cliParse(ctx context.Context, params ...interface{}) (context.Context, error) {
	if len(params) != 1 {
		return ctx, fmt.Errorf("not enough arguments")
	}
	c := params[0].(*cli.Context)
	return context.WithValue(ctx, mainager.Key("github.com/moul/mainager/module/server/grpc.settings.bind-address"), c.String("grpc-bind")), nil
}

func serverInit(ctx context.Context, params ...interface{}) (context.Context, error) {
	options := []grpc.ServerOption{}
	server := grpc.NewServer(options...)

	return context.WithValue(ctx, mainager.Key("github.com/moul/mainager/module/server/grpc.server"), server), nil
}

func serverStart(ctx context.Context, params ...interface{}) (context.Context, error) {
	if len(params) != 1 {
		return ctx, fmt.Errorf("not enough arguments")
	}
	errc := params[0].(chan error)

	server := ctx.Value(mainager.Key("github.com/moul/mainager/module/server/grpc.server")).(*grpc.Server)
	address := ctx.Value(mainager.Key("github.com/moul/mainager/module/server/grpc.settings.bind-address")).(string)
	if address == "" {
		return ctx, nil
	}

	log.Printf("mainager.module.server.grpc listen: %q", address)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return ctx, err
	}

	go func() {
		errc <- server.Serve(ln)
	}()

	return ctx, nil
}
