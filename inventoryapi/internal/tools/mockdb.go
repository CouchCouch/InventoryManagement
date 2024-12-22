package tools

import (
	"inventoryapi/api"
	"time"
)

type mockDB struct{}

var mockItems = []api.Item{
    {
        Id: 2,
        Name: "crash pad",
        Description: "climbing crash pad",
        Quantity: 2,
    },
    {
        Id: 1,
        Name: "microspikes",
        Description: "",
        Quantity: 1,
    },
}

func (d *mockDB) GetItems() *[]api.Item {
    time.Sleep(time.Second * 1)

    var clientData = mockItems

    return &clientData
}

func (d *mockDB) SetupDatabase() error {
    return nil
}
