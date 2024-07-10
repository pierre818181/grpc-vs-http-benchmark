package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type randomStruct struct {
	Message string `json:"message"`
}

func main() {
	ctx := context.Background()
	logger := zaplogger()

	srv := server(ctx)
	client(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server shutdown gracefully")
}

func zaplogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := cfg.Build()
	defer logger.Sync()

	return logger
}

func server(ctx context.Context) *http.Server {
	logger := zaplogger()

	r := gin.New()

	r.GET("/", func(c *gin.Context) {
		var req randomStruct
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response := randomStruct{Message: "Hello " + req.Message}
		c.JSON(http.StatusOK, response)
	})
	r.POST("/", func(c *gin.Context) {
		var req randomStruct
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response := randomStruct{Message: "Hello " + req.Message}
		c.JSON(http.StatusOK, response)
	})

	srv := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	go func() {
		fmt.Println("Starting Gin server on :8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	time.Sleep(1 * time.Second)
	return srv
}

func client(ctx context.Context) {
	logger := zaplogger()
	numOfRequests := 10000

	start := time.Now()
	for i := 0; i < numOfRequests; i++ {
		url := "http://localhost:8081"
		_, err := http.Get(url)
		if err != nil {
			logger.Error("something went wrong", zap.Error(err))
			continue
		}
	}
	end := time.Now()
	logger.Info(fmt.Sprintf("Http time taken for %d GET requests: %v", numOfRequests, end.Sub(start)))

	start = time.Now()
	for i := 0; i < numOfRequests; i++ {
		url := "http://localhost:8081"
		a := randomStruct{
			Message: "Hello world",
		}
		reqBody, _ := json.Marshal(a)
		_, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			logger.Error("something went wrong", zap.Error(err))
			continue
		}
	}

	end = time.Now()
	logger.Info(fmt.Sprintf("Http time taken for %d POST requests: %v", numOfRequests, end.Sub(start)))
}
