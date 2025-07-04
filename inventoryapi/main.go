package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"inventoryapi/internal/handlers"
)

func main() {
	log.SetReportCaller(true)
	var r *chi.Mux = chi.NewRouter()
	handlers.Handler(r)

	fmt.Println("Starting GO API service . . .")
	fmt.Println(`
  _____                      _                              _____ _____
 |_   _|                    | |                       /\   |  __ \_   _|
   | |  _ ____   _____ _ __ | |_ ___  _ __ _   _     /  \  | |__) || |
   | | | '_ \ \ / / _ \ '_ \| __/ _ \| '__| | | |   / /\ \ |  ___/ | |
  _| |_| | | \ V /  __/ | | | || (_) | |  | |_| |  / ____ \| |    _| |_
 |_____|_| |_|\_/ \___|_| |_|\__\___/|_|   \__, | /_/    \_\_|   |_____|
                                            __/ |
                                           |___/`)
	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		log.Error(err)
	}
}
