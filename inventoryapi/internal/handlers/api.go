package handlers

import (
    "github.com/go-chi/chi/v5"
    chimiddle "github.com/go-chi/chi/v5/middleware"
    "inventoryapi/internal/middleware"
)

func Handler(r *chi.Mux) {
    r.Use(chimiddle.StripSlashes)

    r.Route("/user", func(router chi.Router) {

        router.Use(middleware.Authorization)

        router.Get("/items", GetItems)
    })
}
