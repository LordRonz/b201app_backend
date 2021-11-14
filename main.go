package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	"github.com/lordronz/b201app_backend/config"
	"github.com/lordronz/b201app_backend/docs"
	"github.com/lordronz/b201app_backend/pkg/api"
	"github.com/lordronz/b201app_backend/pkg/db"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var addr string
var postgresDSN string

func init() {
	if err := config.LoadConfig("."); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	addr = viper.GetString("address")

	dbUser := viper.GetString("db_user")
	dbHost := viper.GetString("db_host")
	dbPassword := viper.GetString("db_password")
	dbName := viper.GetString("db_name")
	dbPort := viper.GetString("db_port")
	dbSSLMode := viper.GetString("db_sslmode")
	dbTimezone := viper.GetString("db_timezone")

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
