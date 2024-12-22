package handlers

import (
    "github.com/go-chi/chi/v5"
    chimiddle "github.com/go-chi/chi/v5/middleware"
)

func Handler(r *chi.Mux) {
    r.Use(chimiddle.StripSlashes)
    r.Use(chimiddle.Logger)

    r.Get("/items", GetItems)
}
