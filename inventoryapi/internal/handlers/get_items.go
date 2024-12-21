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

    var itemDetails *tools.UserItemDetails
    itemDetails = (*database).GetUserItems(params.Username)
    if(itemDetails == nil) {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }

    var response = api.ItemResponse{
        Code: http.StatusOK,
        Item: api.Item{
            Id: (itemDetails).Id,
            Name: (*itemDetails).Name,
            Description: (*itemDetails).Description,
            Quantity: (*itemDetails).Quantity,
            CheckoutDate: (*itemDetails).CheckoutDate,
        },
    }

    w.Header().Set("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
}
