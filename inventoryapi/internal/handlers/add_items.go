package handlers

import (
	"encoding/json"
	"net/http"

	"inventoryapi/api"
	"inventoryapi/internal/tools"

	// "github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

func AddItems(w http.ResponseWriter, r *http.Request) {
    var params = api.NewItem{}
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

    var itemId *int
    itemId = (*database).AddItem(params)
    if(itemId == nil) {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }

    var response = api.NewItemResponse{
        Code: http.StatusOK,
        Id: *itemId,
    }

    api.EnableCors(&w)

    w.Header().Add("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
}
