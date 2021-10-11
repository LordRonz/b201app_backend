package types

import "net/http"

type PutUserValidation struct {
	Name  string `validate:"required,min=1,max=100" json:"name"`
	Email string `validate:"required,min=1,max=200" json:"email"`
	NRP   string `validate:"required,min=1,max=30" json:"nrp"`
}

func (u *PutUserValidation) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Bind implements the the github.com/go-chi/render.Binder interface
func (u *PutUserValidation) Bind(r *http.Request) error {
	return nil
}
