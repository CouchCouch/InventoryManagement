package tools

import (
	log "github.com/sirupsen/logrus"
    "inventoryapi/api"
)

type DatabaseInterface interface {
    GetItems() *[]api.Item
    GetItem(int) *[]api.Item
    AddItem(api.NewItem) *int
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
