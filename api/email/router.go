package email

import "github.com/go-chi/chi/v5"

func AddRoutes(r chi.Router) {
	r.Post("/", GetEmails)
}
