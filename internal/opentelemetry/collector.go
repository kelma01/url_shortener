package opentelemetry

import (
	"context"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/logging/logrus"
	"gorm.io/plugin/opentelemetry/tracing"
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
	db, err := gorm.Open(sqlite.Open("file:memory:?cache?shared"), &gorm.Config{Logger: logger})
	if err != nil {
		panic(err)
	}
	if db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		panic(err)
	}
	log.Println("Opentelemetry OK")
}

func InitTracer() func() {
	output_file := "traces.log"
	f, err := os.OpenFile(output_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Failed to open traces.log: " + err.Error())
	}
	exporter, _ := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),	//terminale basilmasi icin
		//stdouttrace.WithWriter(f),		//disari file'a export edilmesi icin
	)
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(tp)
	return func() { 
		_ = tp.Shutdown(context.Background())
		_ = f.Close()
	}
}