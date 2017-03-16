package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/moul/mainager"
	_ "github.com/moul/mainager/module/backend/raven"
	_ "github.com/moul/mainager/module/backend/smtp"
	_ "github.com/moul/mainager/module/server/http"
)

func main() {
	var (
		ctx = context.Background()
		m   = mainager.Instance()
		err error
	)

	// Custom Settings
	ctx = context.WithValue(ctx, mainager.Key("github.com/moul/mainager/module/server/http.settings.bind-address"), ":8080")
	ctx = context.WithValue(ctx, mainager.Key("github.com/moul/mainager/module/backend/raven.settings.dsn"), "http://blah:blah@0.0.0.0:0/raven")
	ctx = context.WithValue(ctx, mainager.Key("github.com/moul/mainager/module/backend/smtp.settings.url"), "smtp://user:pass@host:port")

	// Init mainager modules
	if ctx, err = m.InvokeAll(ctx, "backend-init"); err != nil {
		panic(err)
	}
	if ctx, err = m.InvokeAll(ctx, "server-init"); err != nil {
		panic(err)
	}

	// Custom HTTP handlers
	mux := ctx.Value(mainager.Key("github.com/moul/mainager/module/server/http.mux")).(*http.ServeMux)
	mux.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "pong\n")
	})
	mux.HandleFunc("/raven", func(w http.ResponseWriter, req *http.Request) {
		raven := ctx.Value(mainager.Key("github.com/moul/mainager/module/backend/raven.client"))
		fmt.Fprintf(w, "raven: %v\n", raven)
	})
	mux.HandleFunc("/smtp", func(w http.ResponseWriter, req *http.Request) {
		smtp := ctx.Value(mainager.Key("github.com/moul/mainager/module/backend/smtp.client"))
		fmt.Fprintf(w, "smtp: %v\n", smtp)
	})

	// Start server
	errc := make(chan error, 1)
	if ctx, err = m.InvokeAll(ctx, "server-start", errc); err != nil {
		panic(err)
	}
	log.Fatal(<-errc)
}
