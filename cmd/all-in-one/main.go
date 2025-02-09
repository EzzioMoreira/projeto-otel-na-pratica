// Copyright Dose de Telemetria GmbH
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/app"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/config"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/pkg/telemetry"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

func main() {
	configFlag := flag.String("config", "config.yaml", "path to the config file")
	otelConfigFlag := flag.String("otelconfig", "otel.yaml", "path to the otel config file")

	flag.Parse()

	closer, err := telemetry.Setup(context.Background(), *otelConfigFlag)
	if err != nil {
		fmt.Println("failed to setup telemetry", err)
		return
	}
	defer closer(context.Background())

	ctx, span := otel.Tracer("all-in-one").Start(context.Background(), "main")

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		otelzap.NewCore("all-in-one", otelzap.WithLoggerProvider(global.GetLoggerProvider())),
	)

	logger := zap.New(core)

	logger.Info("starting all-in-one")
	span.AddEvent("starting all-in-one")

	c, err := config.LoadConfig(*configFlag)
	if err != nil {
		span.AddEvent("failed to load config", trace.WithAttributes(attribute.String("error", err.Error())))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Error("failed to load config", zap.Error(err))
		return
	}

	mux := http.NewServeMux()

	// starts the gRPC server
	lis, err := net.Listen("tcp", c.Server.Endpoint.GRPC)
	if err != nil {
		span.AddEvent("failed to load config", trace.WithAttributes(attribute.String("error", err.Error())))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Error("failed to listen", zap.Error(err))
		return
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	{
		logger.Info("starting user service")
		span.AddEvent("starting user service")
		a := app.NewUser(ctx, &c.Users)
		a.RegisterRoutes(mux)
	}

	{
		logger.Info("starting plan service")
		span.AddEvent("starting plan service")
		a := app.NewPlan(ctx, &c.Plans)
		a.RegisterRoutes(mux, grpcServer)
	}

	{
		span.AddEvent("starting payment service")
		logger.Info("starting payment service")
		a, err := app.NewPayment(&c.Payments)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			logger.Error("failed to create payment service", zap.Error(err))
		}
		a.RegisterRoutes(mux)
	}

	{
		logger.Info("starting subscription service")
		span.AddEvent("starting subscription service")
		a := app.NewSubscription(&c.Subscriptions)
		a.RegisterRoutes(mux)
	}

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			span.AddEvent("failed to load config", trace.WithAttributes(attribute.String("error", err.Error())))
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			logger.Error("failed to serve", zap.Error(err))
		}
	}()

	span.End()

	err = http.ListenAndServe(c.Server.Endpoint.HTTP, mux)
	if err != nil && err != http.ErrServerClosed {
		span.AddEvent("failed to load config", trace.WithAttributes(attribute.String("error", err.Error())))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Error("failed to serve", zap.Error(err))
	}

	logger.Info("stopping all-in-one")
}
