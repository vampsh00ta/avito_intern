package main

import (
	"avito/config"
	db "avito/internal/db"
	"avito/internal/service"
	"avito/internal/transport"
	postgresql "avito/pkg/client"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type person struct {
	Name     string
	LastName string
	Age      uint8
}

func main() {
	//config
	cfg := config.Load()
	//logger
	logger := LoadLoggerDev()

	//
	dbClient, err := postgresql.NewClient(context.Background(), 5, cfg.DBConfig)
	if err != nil {
		logger.Fatal(err.Error())
	}
	//postgres repository
	repository := db.New(dbClient, logger)

	//service
	srvc := service.New(repository)

	//transport
	httpServer := transport.NewHttpServer(srvc, logger)

	router := mux.NewRouter()
	fmt.Println(router, httpServer)
	server := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}
	router.Methods("POST").Path("/segments").HandlerFunc(httpServer.AddSegment)
	router.Methods("DELETE").Path("/segments").HandlerFunc(httpServer.DeleteSegment)

	router.Methods("POST").Path("/user/new").HandlerFunc(httpServer.CreateUser)
	router.Methods("DELETE").Path("/user").HandlerFunc(httpServer.DeleteUser)

	router.Methods("GET").Path("/user/segments/{id}/").HandlerFunc(httpServer.GetUsersSegments)
	router.Methods("POST").Path("/user/segments").HandlerFunc(httpServer.GetUsersSegments)
	router.Methods("DELETE").Path("/user/segments").HandlerFunc(httpServer.GetUsersSegments)

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
	log.Printf("Starting HTTP server on %s", cfg.Address)

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
