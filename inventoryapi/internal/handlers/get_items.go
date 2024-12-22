package handlers

import (
    "encoding/json"
    "net/http"

    "inventoryapi/api"
    "inventoryapi/internal/tools"
    log "github.com/sirupsen/logrus"
    "github.com/gorilla/schema"
)

func GetItems(w http.ResponseWriter, r *http.Request) {
    var params = api.UserItemParams{}
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

    var items *[]api.Item
    items = (*database).GetItems()
    if(items == nil) {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }

    var response = api.ItemResponse{
        Code: http.StatusOK,
        Items: *items,
    }

    w.Header().Set("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
}
