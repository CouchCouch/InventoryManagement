package handlers

import (
	"encoding/json"
	"net/http"

	"inventoryapi/api"
	"inventoryapi/internal/tools"

	// "github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

func CheckoutItem(w http.ResponseWriter, r *http.Request) {
	params := api.CheckoutParams{}
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

	receipt, err := (*database).CheckoutItem(params)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	response := api.CheckoutItemResponse{
		Code:    200,
		Receipt: *receipt,
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}
