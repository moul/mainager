package sqlite_backend

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/moul/mainager"
	"github.com/urfave/cli"
)

func init() {
	mainager.Register(mainager.Module{
		Name: "github.com/moul/mainager/module/backend/sqlite",
		Hooks: mainager.Hooks{
			"cli-init":     cliInit,
			"cli-parse":    cliParse,
			"backend-init": backendInit,
			"backend-stop": backendStop,
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
			Name:   "sqlite-db-path",
			Usage:  "SQLite database path. If empty, sqlite will be disabled",
			EnvVar: "SQLITE_DB_PATH",
			Value:  "/tmp/sqlite.db",
		},
	}...)

	return ctx, nil
}

func cliParse(ctx context.Context, params ...interface{}) (context.Context, error) {
	if len(params) != 1 {
		return ctx, fmt.Errorf("not enough arguments")
	}
	c := params[0].(*cli.Context)

	return context.WithValue(ctx, mainager.Key("github.com/moul/mainager/module/backend/sqlite.settings.db-path"), c.String("sqlite-db-path")), nil
}

func backendInit(ctx context.Context, params ...interface{}) (context.Context, error) {
	dbPath := ctx.Value(mainager.Key("github.com/moul/mainager/module/backend/sqlite.settings.db-path")).(string)
	if dbPath == "" {
		return ctx, nil
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, mainager.Key("github.com/moul/mainager/module/backend/sqlite.db"), db), nil
}

func backendStop(ctx context.Context, params ...interface{}) (context.Context, error) {
	db := ctx.Value(mainager.Key("github.com/moul/mainager/module/backend/sqlite.db")).(*sql.DB)
	if db != nil {
		db.Close()
	}
	return ctx, nil
}
