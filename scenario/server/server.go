package server_scenario

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/moul/mainager"
	"github.com/urfave/cli"

	_ "github.com/moul/mainager/module/server/grpc"
	_ "github.com/moul/mainager/module/server/http"
)

func Run(ctx context.Context) error {
	var (
		m   = mainager.Instance()
		err error
	)

	app := cli.NewApp()
	app.Name = os.Args[0]
	app.Flags = []cli.Flag{}

	if ctx, err = m.InvokeAll(ctx, "cli-init", app); err != nil {
		panic(err)
	}

	app.Action = func(c *cli.Context) error {
		if ctx, err = m.InvokeAll(ctx, "cli-parse", c); err != nil {
			panic(err)
		}
		if ctx, err = m.InvokeAll(ctx, "backend-init"); err != nil {
			panic(err)
		}
		defer m.InvokeAll(ctx, "backend-stop")

		if ctx, err = m.InvokeAll(ctx, "server-init"); err != nil {
			panic(err)
		}
		if ctx, err = m.InvokeAll(ctx, "service-init"); err != nil {
			panic(err)
		}

		errc := make(chan error, 1)
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			errc <- fmt.Errorf("%s", <-c)
		}()
		if ctx, err = m.InvokeAll(ctx, "server-start", errc); err != nil {
			panic(err)
		}
		return <-errc
	}

	return app.Run(os.Args)
}
