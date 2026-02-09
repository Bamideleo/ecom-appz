package main

import (
	"context"
	"ecom-appz/internal/config"
	"ecom-appz/internal/db"
	"ecom-appz/internal/logger"
	_ "ecom-appz/internal/logger"
	_ "ecom-appz/internal/middleware"
	"ecom-appz/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.SetupEnv()
	 logger := logger.New()
	
	if err != nil{
		log.Fatalf("config file is not loaded properly %v\n", err)
	}

	handler := router.New(logger)

	server := &http.Server{
		Addr:    cfg.AppPort,
		Handler: handler,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error: " + err.Error())
		}
	}()

	<-ctx.Done()
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("forced shutdown: " + err.Error())
	}

	logger.Info("server exited cleanly")


	db.StartServer(cfg)
}