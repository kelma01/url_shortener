package opentelemetry

import (
	"time"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/logging/logrus"
	"gorm.io/plugin/opentelemetry/tracing"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"context"
)

func init() {
	logger := logger.New(
		logrus.NewWriter(),
		logger.Config{
			SlowThreshold: time.Millisecond,
			LogLevel:      logger.Warn,
			Colorful:      false,
		},
	)
	_, _ = gorm.Open(sqlite.Open("file:memory:?cache?shared"), &gorm.Config{Logger: logger})
}

func InitTracing(db *gorm.DB) error {
	return db.Use(tracing.NewPlugin(tracing.WithoutMetrics()))
}

func InitTracer() func() {
	exporter, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(tp)
	return func() { _ = tp.Shutdown(context.Background())}
}