package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	_ "github.com/joho/godotenv/autoload"
	m "github.com/lordronz/b201app_backend/pkg/middleware"
	"github.com/lordronz/b201app_backend/pkg/types"
	"github.com/lordronz/b201app_backend/pkg/db"
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
	// gracefully exit on keyboard interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)

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

	r.Use(m.SetLogger(log))
	r.Use(middleware.Recoverer)
	r.Use(corsConfig().Handler)
	user := &types.User{
		ID:    "1",
		Name:  "Username",
		Email: "email@mail.com",
	}
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, user)
	})
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

func corsConfig() *cors.Cors {
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	return cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           86400, // Maximum value not ignored by any of major browsers
	})
}
