package server_scenario

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/moul/mainager"
	"github.com/urfave/cli"
)

func Run() error {
	m := mainager.Instance()

	app := cli.NewApp()
	app.Name = os.Args[0]

	if err := m.InvokeAll("cli-init", app); err != nil {
		panic(err)
	}

	app.Action = action
	return app.Run(os.Args)
}

func action(c *cli.Context) error {
	m := mainager.Instance()

	if err := m.InvokeAll("cli-parse", c); err != nil {
		panic(err)
	}

	if err := m.InvokeAll("backend-init"); err != nil {
		panic(err)
	}

	if err := m.InvokeAll("server-init"); err != nil {
		panic(err)
	}

	if err := m.InvokeAll("service-init"); err != nil {
		panic(err)
	}

	errc := make(chan error, 1)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()
	if err := m.InvokeAll("server-start", errc); err != nil {
		panic(err)
	}

	return <-errc
}
