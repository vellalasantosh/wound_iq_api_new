package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vellalasantosh/wound_iq_api_new/internal/config"
	"github.com/vellalasantosh/wound_iq_api_new/internal/db"
	"github.com/vellalasantosh/wound_iq_api_new/internal/logger"
	"github.com/vellalasantosh/wound_iq_api_new/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("failed to load config: %v", err)
		os.Exit(1)
	}

	log := logger.New(cfg)
	defer log.Sync()

	sqlDB, err := db.Open(cfg.DB_DSN)
	if err != nil {
		log.Sugar().Fatalf("db open failed: %v", err)
	}
	defer sqlDB.Close()

	r := router.New(sqlDB, log, cfg)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	// start server
	go func() {
		log.Sugar().Infof("starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Sugar().Fatalf("listen: %s", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Sugar().Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Sugar().Fatal("server forced to shutdown:", err)
	}
	log.Sugar().Info("server exiting")
}
