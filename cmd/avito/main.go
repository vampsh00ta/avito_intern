package main

import (
	"avito/config"
	db "avito/internal/db"
	r "avito/internal/redis"
	"avito/internal/service"
	"avito/internal/transport"
	"avito/internal/ttl"
	postgresql "avito/pkg/client"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type test struct {
	time time.Time
	slug string
}

func main() {
	//config
	cfg := config.Load()
	//logger
	logger := LoadLoggerDev()

	//postgres client
	dbClient, err := postgresql.NewClient(context.Background(), 5, cfg.DBConfig)
	if err != nil {
		logger.Fatal(err.Error())
	}
	//redis client
	clientRedis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	//postgres repository
	repository := db.New(dbClient, logger)
	//redis repository
	repRedis := r.New(clientRedis, logger)
	//new TTL
	//service
	ttl := ttl.NewTTL(repository, logger, repRedis)

	srvc := service.New(repository, ttl)

	//transport
	httpServer := transport.NewHttpServer(srvc, logger)

	router := mux.NewRouter()
	server := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	router.Methods("POST").Path("/segments").HandlerFunc(httpServer.AddSegment)
	router.Methods("DELETE").Path("/segments").HandlerFunc(httpServer.DeleteSegment)

	router.Methods("POST").Path("/user/new").HandlerFunc(httpServer.CreateUser)
	router.Methods("DELETE").Path("/user").HandlerFunc(httpServer.DeleteUser)

	router.Methods("GET").Path("/user/segments/{id}/").HandlerFunc(httpServer.GetUsersSegments)
	router.Methods("POST").Path("/user/segments").HandlerFunc(httpServer.AddSegmentsToUser)
	router.Methods("DELETE").Path("/user/segments").HandlerFunc(httpServer.DeleteSegmentsFromUser)

	router.Methods("GET").Path("/history").HandlerFunc(httpServer.GetHistory)

	exit := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(exit)
	}()

	logger.Infow(fmt.Sprintf("Starting HTTP server on %s", cfg.Address))
	go ttl.Start(context.Background(), exit)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-exit
	os.Exit(0)
}
func LoadLoggerDev() *zap.SugaredLogger {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	sugar := logger.Sugar()
	return sugar
}
