package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devopsdemo/order-service/internal/config"
	"github.com/devopsdemo/order-service/internal/db"
	"github.com/devopsdemo/order-service/internal/handlers"
	"github.com/devopsdemo/order-service/internal/kafkaproducer"
	"github.com/gin-gonic/gin"
)

func main() {
	// Structured JSON logging via Go's standard log/slog — no third-party
	// logging library needed, matches the JSON logging pattern used in the
	// Auth and Notification services.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load() // fails fast (os.Exit) inside if prod env vars are missing

	pool, err := db.Connect(cfg)
	if err != nil {
		slog.Error("db_connect_failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	producer := kafkaproducer.New(cfg)
	defer producer.Close()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(handlers.LoggingMiddleware())

	router.GET("/health/live", handlers.Liveness)
	router.GET("/health/ready", handlers.Readiness(pool))

	orderHandler := handlers.NewOrderHandler(pool, producer)
	router.POST("/orders", orderHandler.CreateOrder)
	router.GET("/orders/:id", orderHandler.GetOrder)
	router.GET("/orders", orderHandler.ListOrders)
	router.PATCH("/orders/:id/status", orderHandler.UpdateStatus)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Run the server in a goroutine so main() stays free to listen for the
	// shutdown signal below.
	go func() {
		slog.Info("server_starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server_failed", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown: same intent as the manual SIGTERM/SIGINT handling
	// in the Auth (Node.js) service and server.shutdown=graceful in the
	// Notification (Spring Boot) service. On SIGTERM (what K8s sends before
	// killing a pod), stop accepting new requests but let in-flight ones
	// finish, up to the timeout below.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutdown_signal_received")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("graceful_shutdown_failed", "error", err)
	} else {
		slog.Info("graceful_shutdown_complete")
	}
}
