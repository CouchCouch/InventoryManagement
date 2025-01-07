package handlers

import (
	"encoding/json"
	"net/http"

	"inventoryapi/api"
	"inventoryapi/internal/tools"

	// "github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

func UpdateItem(w http.ResponseWriter, r *http.Request) {
    var params = api.Item{}
    var decoder *json.Decoder = json.NewDecoder(r.Body) // *schema.Decoder = schema.NewDecoder()
    var err error

    err = decoder.Decode(&params)

    log.Printf("%s", params)

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

    var success bool
    success = (*database).UpdateItem(params)
    if success == false {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }

    var response = http.StatusOK

    w.Header().Add("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
}
