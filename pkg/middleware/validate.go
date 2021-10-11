package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"github.com/lordronz/b201app_backend/pkg/types"
)

var validate *validator.Validate = validator.New();

// Pagination middleware is used to extract the next page id from the url query
func Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := &types.PutUserValidation{}
		if err := render.Bind(r, user); err != nil {
			_ = render.Render(w, r, types.ErrInvalidRequest(err))
			return
		}
		if err := validate.Struct(user); err != nil {
			_ = render.Render(w, r, types.ErrInvalidRequest(err))
			return
		}
		ctx := context.WithValue(r.Context(), "put_user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
