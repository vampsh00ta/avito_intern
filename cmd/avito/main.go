package main

import (
	"avito/config"
	_ "avito/docs"
	db "avito/internal/db"
	"avito/internal/service"
	"avito/internal/transport"
	"avito/internal/ttl"
	postgresql "avito/pkg/client"
	"context"
	"fmt"
	"github.com/gorilla/mux"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//query.collection.format multi
//	@host			localhost:8000
//	@BasePath		/api/v1

func main() {

	//logger
	logger := LoadLoggerDev()

	//config
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal(err.Error())
	}
	//postgres client
	dbClient, err := postgresql.NewClient(context.Background(), 5, cfg.DBConfig)
	if err != nil {
		logger.Fatal(err.Error())
	}
	//redis client
	//clientRedis := redis.NewClient(&redis.Options{
	//	Addr:     cfg.Addr,
	//	Password: cfg.Redis.Password,
	//	DB:       cfg.DB,
	//})
	//postgres repository
	repository := db.New(dbClient, logger)
	//redis repository
	//repRedis := r.New(clientRedis, logger)
	//new TTL
	//service
	//ttl := ttl.New(repository, logger, repRedis, cfg)
	ttl := ttl.New(repository, logger, cfg)

	srvc := service.New(repository, ttl)

	//transport
	httpServer := transport.NewHttpServer(srvc, logger)

	router := mux.NewRouter()
	server := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	router.Methods("POST").Path("/user/new").HandlerFunc(httpServer.CreateUser)
	router.Methods("DELETE").Path("/user").HandlerFunc(httpServer.DeleteUser)
	router.Methods("GET").Path("/user/segments/{id}").HandlerFunc(httpServer.GetUsersSegments)
	router.Methods("POST").Path("/user/segments/new").HandlerFunc(httpServer.AddSegmentsToUser)
	router.Methods("DELETE").Path("/user/segments").HandlerFunc(httpServer.DeleteSegmentsFromUser)

	router.Methods("POST").Path("/segment/new").HandlerFunc(httpServer.AddSegment)
	router.Methods("DELETE").Path("/segment").HandlerFunc(httpServer.DeleteSegment)

	router.Methods("GET").Path("/history").HandlerFunc(httpServer.GetHistory)
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

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
	go ttl.Start(exit)

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
