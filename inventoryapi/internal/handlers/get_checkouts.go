package handlers

import (
	"encoding/json"
	"net/http"

	"inventoryapi/api"
	"inventoryapi/internal/tools"

	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

func GetCheckouts(w http.ResponseWriter, r *http.Request) {
    var params = api.ItemParams{}
    var decoder = schema.NewDecoder()
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

    var checkouts *[]api.CheckoutItem
    if params.Id != 0 {
        checkouts, err = (*database).GetCheckout(params.Id)
    } else {
        checkouts, err = (*database).GetCheckouts()
    }
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }

    var response = api.CheckoutResponse{
        Code: http.StatusOK,
        Checkouts: *checkouts,
    }

    w.Header().Add("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
}
