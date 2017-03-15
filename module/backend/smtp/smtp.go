package raven_backend

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/moul/mainager"
	"github.com/urfave/cli"
	gomail "gopkg.in/gomail.v2"
)

func init() {
	mainager.Register(mainager.Module{
		Name: "mainager.module.backend.smtp",
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
			Name:   "smtp-url",
			Usage:  "SMTP server url (i.e, `smtp://user:pass@host:port`). If empty, SMTP will be disabled.",
			EnvVar: "SMTP_URL",
			Value:  "smtp://127.0.0.1:25",
		},
	}...)
	return ctx, nil
}

func cliParse(ctx context.Context, params ...interface{}) (context.Context, error) {
	if len(params) != 1 {
		return ctx, fmt.Errorf("not enough arguments")
	}
	c := params[0].(*cli.Context)
	return context.WithValue(ctx, mainager.Key("mainager.module.backend.smtp.settings.url"), c.String("smtp-url")), nil
}

func backendInit(ctx context.Context, params ...interface{}) (context.Context, error) {
	smtpURL := ctx.Value(mainager.Key("mainager.module.backend.smtp.settings.url")).(string)
	if smtpURL == "" {
		return ctx, nil
	}

	u, err := url.Parse(smtpURL)
	if err != nil {
		return ctx, nil
	}

	if u.Scheme != "smtp" {
		return ctx, fmt.Errorf("unsupported scheme: %q", u.Scheme)
	}
	host, portStr, _ := net.SplitHostPort(u.Host)
	if portStr == "" {
		portStr = "25" // default port value
	}
	port, _ := strconv.ParseInt(portStr, 10, 64)
	var username, password string
	if u.User != nil {
		username = u.User.Username()
		password, _ = u.User.Password()
	}

	client := gomail.NewDialer(host, int(port), username, password)
	ctx = context.WithValue(ctx, mainager.Key("mainager.module.backend.smtp.client"), client)

	return ctx, nil
}
