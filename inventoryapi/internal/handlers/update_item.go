package handlers

import (
	"encoding/json"
	"inventoryapi/api"
	"inventoryapi/internal/tools"
	"net/http"

	// "github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	params := api.Item{}
	var decoder *json.Decoder = json.NewDecoder(r.Body) // *schema.Decoder = schema.NewDecoder()
	var err error

	err = decoder.Decode(&params)

	log.Printf("%s", params.String())

	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	var database *tools.DatabaseInterface
	database, err = tools.NewDatabase()
	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	defer (*database).CloseDatabase()

	err = (*database).UpdateItem(params)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	response := http.StatusOK

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}
