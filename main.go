package main

import (
	"context"
	"errors"
	"fmt"
	"mrrancy/logAssignment/logger"
	"mrrancy/logAssignment/store"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mrrancy/logAssignment/handlers"
	"mrrancy/logAssignment/middlewares"
	"mrrancy/logAssignment/utils"
)

func main() {

	// Create a Logger instance
	log := setupLogger()
	log.Info("Initializing Logger...")

	// Environment Variables
	batchSize := utils.GetEnvAsInt("BATCH_SIZE", 5)
	batchInterval := utils.GetEnvAsInt("BATCH_INTERVAL", 1)
	postEndpoint := utils.GetEnv("POST_ENDPOINT", "https://httpdump.app/dumps/8e06d9f6-ac86-4629-a13b-f2d802a24fbe")

	// Creating In Memory store
	cache := store.NewCache(log, postEndpoint, batchInterval, batchSize, 3, 2)
	cache.InitializeLogMonitor()

	// Initializing the Controller
	controller := handlers.NewController(log, cache)

	r := setupGinAndRoutes(controller)

	serverAddr := fmt.Sprintf(":%s", utils.GetEnv("PORT", "9001"))
	srv := setupServer(serverAddr, r)

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("listen: " + err.Error())
		}
	}()

	signal.Notify(cache.Exit, syscall.SIGINT, syscall.SIGTERM)
	// Wait for termination signal
	<-cache.Exit
	log.Info("Received Exit Signal")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shut down the server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:" + err.Error())
	}
}

func setupLogger() *zap.Logger {
	return logger.Initialize()
}

func setupGinAndRoutes(controller *handlers.Controller) *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/healthz", controller.Health)
	r.POST("/log", middlewares.LogMiddleware(controller.Log), controller.LogHandler)
	return r
}

func setupServer(serverAddr string, r *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}
}
