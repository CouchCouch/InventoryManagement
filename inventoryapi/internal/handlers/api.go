package handlers

import (
	"github.com/go-chi/chi/v5"
	chimiddle "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"
)

func Handler(r *chi.Mux) {
    r.Use(chimiddle.StripSlashes)
    r.Use(chimiddle.Logger)
    r.Use(chimiddle.RequestID)
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"*"}, // Use this to allow specific origin hosts
        // AllowedOrigins:   []string{"https://*", "http://*"},
        // AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300, // Maximum value not ignored by any of major browsers
    }))

    r.Get("/items", GetItems)
    r.Post("/items", AddItems)
    r.Delete("/items", DeleteItems)
    r.Put("/items", UpdateItem)
    r.Post("/checkout", CheckoutItem)
    r.Put("/checkout", ReturnItem)
    r.Get("/checkout", GetCheckouts)
}
