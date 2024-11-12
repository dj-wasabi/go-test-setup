package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"werner-dijkerman.nl/test-setup/internal/adapter/in/http/api"
	"werner-dijkerman.nl/test-setup/internal/adapter/out/mongodb"
	"werner-dijkerman.nl/test-setup/internal/core/domain/services"
	"werner-dijkerman.nl/test-setup/pkg/config"
	"werner-dijkerman.nl/test-setup/pkg/logging"
)

func main() {
	logger := logging.Initialize()
	logger.Info("Starting the service")

	c := config.ReadConfig()

	con := mongodb.NewMongodbConnection(c)
	repoUser := mongodb.NewUserMongoRepo(con, "users")
	repoOrg := mongodb.NewOrganisationMongoRepo(con, "organisations")

	serviceUser := mongodb.NewUserMongoService(repoUser, logger)
	serviceOrganisation := mongodb.NewOrganisationMongoService(repoOrg, logger)

	h := services.NewdomainServices(serviceOrganisation, serviceUser)

	server := api.NewGinServer(serviceUser, api.NewApiService(h), c, logger)
	server.ListenAndServe()

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
