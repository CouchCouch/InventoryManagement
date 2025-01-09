package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type ItemParams struct {
    Id int
}

type Item struct {
    Id int
    Name string
    Description string
    Quantity int
}

type ItemResponse struct {
    Code int
    Items []Item
}

type NewItem struct {
    Name string //`schema:"name,required"`
    Description string //`schema:"description,required"`
    Quantity int //`schema:"quantity,default:1"`
}

type NewItemResponse struct {
    Code int
    Id int
}

type CheckoutItem struct {
    Id int
    Name string
    Email string
}

type CheckoutItemReceipt struct {
    ItemId int
    Name string
    Email string
    Date time.Time
}

type CheckoutItemResponse struct {
    Code int
    Receipt CheckoutItemReceipt
}

type Error struct {
    Code int

    Message string
}

func writeError(w http.ResponseWriter, message string, code int) {
    resp := Error{
        Code: code,
        Message: message,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)

    json.NewEncoder(w).Encode(resp)
}

var (
    RequestErrorHandler = func(w http.ResponseWriter, err error) {
        writeError(w, err.Error(), http.StatusBadRequest)
    }
    InternalErrorHandler = func(w http.ResponseWriter) {
        writeError(w, "An unexpected error has occured", http.StatusInternalServerError)
    }
)
