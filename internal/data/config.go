package data

import (
	"time"

	"github.com/gookit/config/v2"
)

type ConfigAppLogger struct {
	Pretty bool `mapstructure:"pretty"`
}

type ConfigApp struct {
	Name           string          `mapstructure:"name"`
	RedirectorHost string          `mapstructure:"redirector_host"`
	Key            uint32          `mapstructure:"key"`
	Logger         ConfigAppLogger `mapstructure:"logger"`
}

type ConfigHttpServer struct {
	Addr    string `mapstructure:"addr"`
	Prefork bool   `mapstructure:"prefork"`
}

type ConfigDatabaseRedis struct {
	Addr         string        `mapstructure:"addr"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type ConfigDatabase struct {
	Redis ConfigDatabaseRedis `mapstructure:"redis"`
}

type ConfigOpenTelemetry struct {
	TracerEndpoint string `mapstructure:"tracer_endpoint"`
}

type Config struct {
	App        ConfigApp           `mapstructure:"app"`
	HttpServer ConfigHttpServer    `mapstructure:"http_server"`
	Database   ConfigDatabase      `mapstructure:"database"`
	Otel       ConfigOpenTelemetry `mapstructure:"otel"`
}

func NewConfig(global string, app string) (cfg Config, err error) {
	config.WithOptions(config.ParseEnv, config.ParseTime)

	if err = config.LoadFiles(global, app); err != nil {
		return
	}
	if err = config.BindStruct("", &cfg); err != nil {
		return
	}

	return
}
