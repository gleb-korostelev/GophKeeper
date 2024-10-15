package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gleb-korostelev/GophKeeper/internal/config"
	"github.com/gleb-korostelev/GophKeeper/internal/db/dbimpl"
	"github.com/gleb-korostelev/GophKeeper/internal/service/handler"
	"github.com/gleb-korostelev/GophKeeper/internal/service/router"
	"github.com/gleb-korostelev/GophKeeper/internal/storage/repository"
	"github.com/gleb-korostelev/GophKeeper/internal/workerpool"
	logger "github.com/gleb-korostelev/GophKeeper/tools"
	"go.uber.org/zap"
)

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		logger.Infof("Error in logger: %v", err)
		return
	}

	err = config.ConfigInit()
	if err != nil {
		logger.Infof("Error in config: %v", err)
	}
	database, err := dbimpl.InitDB()
	if err != nil {
		logger.Infof("Error database initialize: %v", err)
		return
	}

	store := repository.NewDBStorage(database)
	defer store.Close()

	workerPool := workerpool.NewDBWorkerPool(config.MaxRoutine)
	defer workerPool.Shutdown()
	svc := handler.NewAPIService(store, workerPool)
	r := router.RouterInit(svc, log)

	logger.Infof("Server is listening on: %s", config.ServerConfig.ServerAddr)

	srv := &http.Server{
		Addr:    config.ServerConfig.ServerAddr,
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Error starting server: %v", err)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exiting")

}
