package tools

import (
	log "github.com/sirupsen/logrus"
    "inventoryapi/api"
)

type DatabaseInterface interface {
    GetItems() *[]api.Item
    SetupDatabase() error
}

func NewDatabase() (*DatabaseInterface, error) {

    var database DatabaseInterface = &mockDB{}

    var err error = database.SetupDatabase()
    if err != nil {
        log.Error(err)
        return nil, err
    }

    return &database, nil
}
