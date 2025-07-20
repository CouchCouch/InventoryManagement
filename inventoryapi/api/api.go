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

func (i *Item) String() string {
	return fmt.Sprintf("Item{Id: %d, Name: %s, Description: %s, Quantity: %d}", i.Id, i.Name, i.Description, i.Quantity)
}

func (i *ItemParams) String() string {
	return fmt.Sprintf("ItemParams{Id: %d}", i.Id)
}

func (i *ItemResponse) String() string {
	return fmt.Sprintf("ItemResponse{Code: %d, Items: %v}", i.Code, i.Items)
}

func (n *NewItemResponse) String() string {
	return fmt.Sprintf("NewItemResponse{Code: %d, Id: %d}", n.Code, n.Id)
}

func (c *CheckoutParams) String() string {
	return fmt.Sprintf("CheckoutParams{Id: %d, Name: %s, Email: %s}", c.Id, c.Name, c.Email)
}

func (c *CheckoutItemReceipt) String() string {
	return fmt.Sprintf("CheckoutItemReceipt{ItemId: %d, Name: %s, Email: %s, Date: %v}", c.ItemId, c.Name, c.Email, c.Date)
}

func (c *CheckoutItem) String() string {
	return fmt.Sprintf("CheckoutItem{Id: %d, ItemId: %d, ItemName: %s, Name: %s, Email: %s, Date: %v, Returned: %t}", c.Id, c.ItemId, c.ItemName, c.Name, c.Email, c.Date, c.Returned)
}

func (c *CheckoutItemResponse) String() string {
	return fmt.Sprintf("CheckoutItemResponse{Code: %d, Receipt: %v}", c.Code, c.Receipt)
}

func (c *CheckoutResponse) String() string {
	return fmt.Sprintf("CheckoutResponse{Code: %d, Checkouts: %v}", c.Code, c.Checkouts)
}

func (e *Error) String() string {
	return fmt.Sprintf("Error{Code: %d, Message: %s}", e.Code, e.Message)
}

// hello this is going to be a really long sentence so
