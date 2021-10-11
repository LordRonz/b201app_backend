package api

import (
	"net/http"

	"github.com/go-chi/render"

	m "github.com/lordronz/b201app_backend/pkg/middleware"
	"github.com/lordronz/b201app_backend/pkg/types"
)

// GetUser renders the user from the context
// @Summary Get user by id
// @Description GetUser returns a single user by id
// @Tags Users
// @Produce json
// @Param id path string true "user id"
// @Router /users/{id} [get]
// @Success 200 {object} types.User
// @Failure 400 {object} types.ErrResponse
// @Failure 404 {object} types.ErrResponse
func GetUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(m.UserCtxKey).(*types.User)

	if err := render.Render(w, r, user); err != nil {
		_ = render.Render(w, r, types.ErrRender(err))
		return
	}
}

// PutUser writes an user to the database
// @Summary Add an user to the database
// @Description PutUser writes an user to the database
// @Description To write a new user, leave the id empty. To update an existing one, use the id of the user to be updated
// @Tags Users
// @Produce json
// @Router /users [put]
// @Success 200 {object} types.User
// @Failure 400 {object} types.ErrResponse
// @Failure 404 {object} types.ErrResponse
func PutUser(w http.ResponseWriter, r *http.Request) {
	body := r.Context().Value("put_user").(*types.PutUserValidation);
	user := &types.User{Name: body.Name, Email: body.Email, NRP: body.NRP}

	if err := DBClient.SetUser(user); err != nil {
		_ = render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}

	if err := render.Render(w, r, user); err != nil {
		_ = render.Render(w, r, types.ErrRender(err))
		return
	}
}

// ListUsers returns all users in the database
// @Summary List all users
// @Description Get all users stored in the database
// @Tags Users
// @Produce json
// @Param page_id query string false "id of the page to be retrieved"
// @Router /users [get]
// @Success 200 {object} types.UserList
// @Failure 400 {object} types.ErrResponse
// @Failure 404 {object} types.ErrResponse
func ListUsers(w http.ResponseWriter, r *http.Request) {
	pageID := r.Context().Value(m.PageIDKey)
	if err := render.Render(w, r, DBClient.GetUsers(pageID.(int))); err != nil {
		_ = render.Render(w, r, types.ErrRender(err))
		return
	}
}
