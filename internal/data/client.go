package data

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/jr-dragon/dynamic_link/internal/library/logs"
)

type Clients struct {
	// RDB is the connection of Redis
	RDB *redis.Client
	// TracerProvider is the otel tracer provider
	TracerProvider trace.TracerProvider
}

func NewClients(cfg Config) (*Clients, error) {
	var err error
	c := &Clients{}

	ctx := context.Background()

	if err = c.initTracerProvider(ctx, cfg); err != nil {
		slog.ErrorContext(ctx, "init tracer provider failed", logs.Err(err))
	}
	if err = c.initRedisClient(cfg); err != nil {
		slog.ErrorContext(ctx, "init redis client failed", logs.Err(err))
	}

	return c, nil
}

func (c *Clients) newOtelExporter(ctx context.Context, cfg Config) (*otlptrace.Exporter, error) {
	conn, err := grpc.NewClient(cfg.Otel.TracerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}

func (c *Clients) initTracerProvider(ctx context.Context, cfg Config) error {
	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(cfg.App.Name)))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	exporter, err := c.newOtelExporter(ctx, cfg)
	if err != nil {
		return err
	}

	c.TracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
	)
	otel.SetTracerProvider(c.TracerProvider)

	return nil
}

func (c *Clients) initRedisClient(cfg Config) error {
	c.RDB = redis.NewClient(&redis.Options{
		Addr:         cfg.Database.Redis.Addr,
		ReadTimeout:  cfg.Database.Redis.ReadTimeout,
		WriteTimeout: cfg.Database.Redis.WriteTimeout,
	})

	return redisotel.InstrumentTracing(c.RDB)
}
