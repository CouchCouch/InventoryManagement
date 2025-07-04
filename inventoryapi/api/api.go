package api

import (
	"encoding/json"
	"fmt"
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
	Code  int    `json:"code"`
	Items []Item `json:"items"`
}

type NewItem struct {
	Name        string `json:"name"`        //`schema:"name,required"`
	Description string `json:"description"` //`schema:"description,required"`
	Quantity    int    `json:"quantity"`    //`schema:"quantity,default:1"`
}

type NewItemResponse struct {
	Code int `json:"code"`
	Id   int `json:"id"`
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
	Code    int                 `json:"code"`
	Receipt CheckoutItemReceipt `json:"receipt"`
}

type CheckoutResponse struct {
	Code      int            `json:"code"`
	Checkouts []CheckoutItem `json:"checkouts"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, message string, code int) {
	resp := Error{
		Code:    code,
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

func (n *NewItem) String() string {
	return fmt.Sprintf("NewItem{Name: %s, Description: %s, Quantity: %d}", n.Name, n.Description, n.Quantity)
}
