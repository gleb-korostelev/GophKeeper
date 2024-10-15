package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gleb-korostelev/GophKeeper/cmd/initConnection"
	"github.com/gleb-korostelev/GophKeeper/config"
	"github.com/gleb-korostelev/GophKeeper/tools/closer"
	"github.com/gleb-korostelev/GophKeeper/tools/logger"
)

func main() {
	defer func() {
		closer.Wait()
		closer.CloseAll()
	}()

	ctx := context.Background()

	db := initConnection.NewDBConn()

	port := config.GetConfigInt(config.Port)

	logger.Info(fmt.Sprintf("Started server at :%d. Swagger docs stated at %d", port, port+1))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: initConnection.InitImpl(ctx, db, port),
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
