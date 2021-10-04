package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/lordronz/b201app_backend/pkg/db"
	"github.com/lordronz/b201app_backend/pkg/types"
)

type (
	// CustomKey is used to refer to the context key that stores custom values of this api to avoid overwrites
	CustomKey string
)

const (
	// UserCtxKey refers to the context key that stores the user
	UserCtxKey CustomKey = "user"
)

var DBClient db.ClientInterface

func SetDBClient(c db.ClientInterface) {
	DBClient = c
}

// User middleware is used to load an User object from
// the URL parameters passed through as the request. In case
// the User could not be found, we stop here and return a 404.
func User(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *types.User

		if id := chi.URLParam(r, "id"); id != "" {
			intID, err := strconv.Atoi(id)
			if err != nil {
				_ = render.Render(w, r, types.ErrInvalidRequest(err))
				return
			}
			user = DBClient.GetUserByID(intID)
		} else {
			_ = render.Render(w, r, types.ErrNotFound())
			return
		}
		if user == nil {
			_ = render.Render(w, r, types.ErrNotFound())
			return
		}

		ctx := context.WithValue(r.Context(), UserCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
