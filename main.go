package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	"github.com/lordronz/b201app_backend/pkg/api"
	"github.com/lordronz/b201app_backend/pkg/db"
	"github.com/lordronz/b201app_backend/docs"
	"go.uber.org/zap"
)

var addr string
var postgresDSN string

func init() {
	addr = os.Getenv("ADDRESS")

	dbUser := os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	dbTimezone := os.Getenv("DB_TIMEZONE")

	postgresDSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode, dbTimezone)
}

func main() {
	docs.SwaggerInfo.Version = "1.0.0"

	// gracefully exit on keyboard interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// configure logger
	log, _ := zap.NewProduction(zap.WithCaller(false))
	defer func() {
		_ = log.Sync()
	}()

	dbClient := &db.Client{}
	if err := dbClient.Connect(postgresDSN); err != nil {
		log.Error("couldn't connect to database", zap.Error(err))
		os.Exit(1)
	}

	r := api.GetRouter(log, dbClient)

	go func() {
		if err := http.ListenAndServe(addr, r); err != nil {
			log.Error("failed to start server", zap.Error(err))
			os.Exit(1)
		}
	}()

	log.Info("ready to serve requests on " + addr)
	<-c
	log.Info("gracefully shutting down")
	os.Exit(0)
}
