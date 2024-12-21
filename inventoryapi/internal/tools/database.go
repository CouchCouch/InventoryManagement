package tools

import (
	log "github.com/sirupsen/logrus"
    "time"
)

type LoginDetails struct {
    AuthToken string
    Username string
}

type UserItemDetails struct {
    Username string
    Id int
    Name string
    Description string
    Quantity int
    CheckoutDate time.Time
}

type DatabaseInterface interface {
    GetUserLoginDetails(username string) *LoginDetails
    GetUserItems(username string) *UserItemDetails
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
