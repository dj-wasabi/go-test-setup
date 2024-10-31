package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"werner-dijkerman.nl/test-setup/adapter/in/http/api"
	"werner-dijkerman.nl/test-setup/adapter/out/mongodb"
	"werner-dijkerman.nl/test-setup/domain/services"
	"werner-dijkerman.nl/test-setup/pkg/config"
	"werner-dijkerman.nl/test-setup/pkg/logging"
)

func main() {
	logger := logging.Initialize()
	logger.Info("Starting the service")

	c := config.ReadConfig()

	org, usr := mongodb.NewMongoDBConnection(c)
	h := services.NewdomainServices(org, usr)
	server := api.NewGinServer(api.NewApiService(h), c)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("Failed to initialize server: %v\n", err.Error()))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	stop()
	logger.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Server forced to shutdown: %v\n", err.Error()))
	}

	logger.Info("Server exiting")

}
