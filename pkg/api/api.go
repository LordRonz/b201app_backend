package api

import (
	"github.com/lordronz/b201app_backend/pkg/db"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	httpSwagger "github.com/swaggo/http-swagger"

	m "github.com/lordronz/b201app_backend/pkg/middleware"
)

var DBClient db.ClientInterface

func SetDBClient(c db.ClientInterface) {
	DBClient = c
	m.SetDBClient(DBClient)
}

// GetRouter configures a chi router and starts the http server
// @title B201 App API
// @description This API is a sample go-api.
// @description It also does this.
// @contact.name B201Crew
// @contact.email b201crew@gmail.com
// @host example.com
// @BasePath /
func GetRouter(log *zap.Logger, dbClient db.ClientInterface) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	SetDBClient(dbClient)
	if log != nil {
		r.Use(m.SetLogger(log))
	}
	r.Use(middleware.Recoverer)
	r.Use(corsConfig().Handler)
	buildTree(r)

	return r
}

func buildTree(r *chi.Mux) {
	r.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/", http.StatusMovedPermanently)
	})
	r.Get("/swagger*", httpSwagger.Handler())
	r.Route("/users", func(r chi.Router) {
		r.With(m.Pagination).Get("/", ListUsers)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(m.User)
			r.Get("/", GetUser)
		})
		r.With(m.Validate).Put("/", PutUser)
	})
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
