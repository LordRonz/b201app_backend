package types

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type User struct {
	ID        uint      `gorm:"type:SERIAL;PRIMARY_KEY" json:"id" example:"1"`
	Name      string    `gorm:"type:varchar;NOT NULL" json:"name" example:"Amogus"`
	Email     string    `gorm:"type:varchar;NOT NULL" json:"email" example:"amogus@mail.com"`
	NRP       string    `gorm:"type:varchar(20);NOT NULL" json:"nrp"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (u *User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Bind implements the the github.com/go-chi/render.Binder interface
func (u *User) Bind(r *http.Request) error {
	return nil
}

// UserList contains a list of users
type UserList struct {
	// A list of users
	Items []*User `json:"items"`
	// The id to query the next page
	NextPageID uint `json:"nextPageId,omitempty" example:"69"`
} // @name UserList

// Render implements the github.com/go-chi/render.Renderer interface
func (a *UserList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status" example:"Resource not found."`                                         // user-level status message
	AppCode    int64  `json:"code,omitempty" example:"404"`                                                 // application-specific error code
	ErrorText  string `json:"error,omitempty" example:"The requested resource was not found on the server"` // application-level error message, for debugging
} // @name ErrorResponse

// Render implements the github.com/go-chi/render.Renderer interface for ErrResponse
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest returns a structured http response for invalid requests
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

// ErrRender returns a structured http response in case of rendering errors
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

// ErrNotFound returns a structured http response if a resource couln't be found
func ErrNotFound() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "Resource not found.",
	}
}
