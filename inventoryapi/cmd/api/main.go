package main

import (
    "fmt"
    "net/http"

    "github.com/go-chi/chi/v5"
    "inventoryapi/internal/handlers"
    log "github.com/sirupsen/logrus"
)

func main() {
    log.SetReportCaller(true)
    var r *chi.Mux = chi.NewRouter()
    handlers.Handler(r)

    fmt.Println("Starting GO API service . . .")
    fmt.Println(`
___  ________   ___      ___ _______   ________   _________  ________  ________      ___    ___      ________  ________  ___
|\  \|\   ___  \|\  \    /  /|\  ___ \ |\   ___  \|\___   ___\\   __  \|\   __  \    |\  \  /  /|    |\   __  \|\   __  \|\  \
\ \  \ \  \\ \  \ \  \  /  / | \   __/|\ \  \\ \  \|___ \  \_\ \  \|\  \ \  \|\  \   \ \  \/  / /    \ \  \|\  \ \  \|\  \ \  \
 \ \  \ \  \\ \  \ \  \/  / / \ \  \_|/_\ \  \\ \  \   \ \  \ \ \  \\\  \ \   _  _\   \ \    / /      \ \   __  \ \   ____\ \  \
  \ \  \ \  \\ \  \ \    / /   \ \  \_|\ \ \  \\ \  \   \ \  \ \ \  \\\  \ \  \\  \|   \/  /  /        \ \  \ \  \ \  \___|\ \  \
   \ \__\ \__\\ \__\ \__/ /     \ \_______\ \__\\ \__\   \ \__\ \ \_______\ \__\\ _\ __/  / /           \ \__\ \__\ \__\    \ \__\
    \|__|\|__| \|__|\|__|/       \|_______|\|__| \|__|    \|__|  \|_______|\|__|\|__|\___/ /             \|__|\|__|\|__|     \|__|
                                                                                    \|___|/`)
    err := http.ListenAndServe("localhost:8080", r)
    if err != nil {
        log.Error(err)
    }

}

