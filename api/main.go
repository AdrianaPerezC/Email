package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AdrianaPerezC/Email/email"
	"github.com/AdrianaPerezC/Email/shared"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	email.AddRoutes(r)
	handler := shared.Cors(r)

	fmt.Print("Escuchando correctamente el puerto- " + os.Getenv("PORT"))
	fmt.Print("UREL ---->" + os.Getenv("FRONTEND_URL"))
	http.ListenAndServe(os.Getenv("PORT"), handler)
	if err := http.ListenAndServe(os.Getenv("PORT"), r); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
