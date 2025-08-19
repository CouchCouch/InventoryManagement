package tools

import (
	"inventoryapi/api"

	log "github.com/sirupsen/logrus"
)

type DatabaseInterface interface {
	GetItems() (*[]api.Item, error)
	GetItem(int) (*[]api.Item, error)
	AddItem(api.NewItem) (*int, error)
	UpdateItem(api.Item) error
	DeleteItem(int) error
	CheckoutItem(api.CheckoutParams) (*api.CheckoutItemReceipt, error)
	ReturnItem(int) error
	GetCheckouts() (*[]api.CheckoutItem, error)
	GetCheckout(int) (*[]api.CheckoutItem, error)
	SetupDatabase() error
	CloseDatabase() error
}

func NewDatabase() (*DatabaseInterface, error) {
	var database DatabaseInterface = &sqlDB{}

	var err error = database.SetupDatabase()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &database, nil
}
