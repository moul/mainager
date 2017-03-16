package raven_backend

import (
	"context"
	"fmt"

	raven "github.com/getsentry/raven-go"
	"github.com/moul/mainager"
	"github.com/urfave/cli"
)

func init() {
	mainager.Register(mainager.Module{
		Name: "github.com/moul/mainager/module/backend/raven",
		Hooks: mainager.Hooks{
			"cli-init":     cliInit,
			"cli-parse":    cliParse,
			"backend-init": backendInit,
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
			Name:   "raven-dsn",
			Usage:  "Raven DSN. If empty, raven will be disabled",
			EnvVar: "RAVEN_DSN",
		},
	}...)
	return ctx, nil
}

func cliParse(ctx context.Context, params ...interface{}) (context.Context, error) {
	if len(params) != 1 {
		return ctx, fmt.Errorf("not enough arguments")
	}
	c := params[0].(*cli.Context)
	return context.WithValue(ctx, mainager.Key("github.com/moul/mainager/module/backend/raven.settings.dsn"), c.String("raven-dsn")), nil
}

func backendInit(ctx context.Context, params ...interface{}) (context.Context, error) {
	dsn := ctx.Value(mainager.Key("github.com/moul/mainager/module/backend/raven.settings.dsn")).(string)
	if dsn == "" {
		return ctx, nil
	}

	client, err := raven.New(dsn)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, mainager.Key("github.com/moul/mainager/module/backend/raven.client"), client), nil
}
