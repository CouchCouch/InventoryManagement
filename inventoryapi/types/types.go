package types

import (
	"net/http"
	"time"
)

type ItemParams struct {
	Id int `schema:"id"`
}

type Item struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

type ItemResponse struct {
	Items []Item `json:"items"`
}

type NewItem struct {
	Name        string `json:"name"`        //`schema:"name,required"`
	Description string `json:"description"` //`schema:"description,required"`
	Quantity    int    `json:"quantity"`    //`schema:"quantity,default:1"`
}

type NewItemResponse struct {
	Id int `json:"id"`
}

type CheckoutParams struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CheckoutItemReceipt struct {
	ItemId int       `json:"itemId"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Date   time.Time `json:"date"`
}

type CheckoutItem struct {
	Id       int       `json:"id"`
	ItemId   int       `json:"itemId"`
	ItemName string    `json:"itemName"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Date     time.Time `json:"date"`
	Returned bool      `json:"returned"`
}

type CheckoutItemResponse struct {
	Receipt CheckoutItemReceipt `json:"receipt"`
}

type CheckoutResponse struct {
	Checkouts []CheckoutItem `json:"checkouts"`
}

func writeError(w http.ResponseWriter, message string, code int) {
	http.Error(w, message, code)
}

var (
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusBadRequest)
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "An unexpected error has occured", http.StatusInternalServerError)
	}
)
