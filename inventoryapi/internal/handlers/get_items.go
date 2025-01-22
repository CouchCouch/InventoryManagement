package handlers

import (
	"encoding/json"
	"net/http"

	"inventoryapi/api"
	"inventoryapi/internal/tools"

	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

func GetItems(w http.ResponseWriter, r *http.Request) {
    var params = api.ItemParams{}
    var decoder *schema.Decoder = schema.NewDecoder()
    var err error

    err = decoder.Decode(&params, r.URL.Query())

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

    var items *[]api.Item
    if params.Id != 0 {
        items, err = (*database).GetItem(params.Id)
    } else {
        items, err = (*database).GetItems()
    }

    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }

    var response = api.ItemResponse{
        Code: http.StatusOK,
        Items: *items,
    }

    //api.EnableCors(&w)

    w.Header().Add("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
}
