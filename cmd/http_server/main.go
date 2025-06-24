package main

import (
	"flag"
	"github.com/jr-dragon/dynamic_link/internal/data"
	"github.com/jr-dragon/dynamic_link/internal/library/logs"
	"log/slog"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagCfgPath is the config path.
	flagGlobalCfgPath string
	flagAppCfgPath    string
)

func init() {
	flag.StringVar(&flagGlobalCfgPath, "global_config", "config.json", "project shared config")
	flag.StringVar(&flagAppCfgPath, "app_config", "cmd/http_server/config.json", "app config")
}

func main() {
	flag.Parse()

	cfg, err := data.NewConfig(flagAppCfgPath, flagGlobalCfgPath)
	if err != nil {
		slog.Error("failed to load config", logs.Err(err))
	}

	app, cleanup, err := wireApp(cfg)
	if err != nil {
		slog.Error("failed to initialize app", logs.Err(err))
	}
	defer cleanup()

	err = app.Listen(cfg.HttpServer.Addr)
}
