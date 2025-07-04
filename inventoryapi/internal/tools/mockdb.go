package tools

import (
	"time"

	"inventoryapi/api"
)

type mockDB struct{}

var mockItems = []api.Item{
	{
		Id:          2,
		Name:        "crash pad",
		Description: "climbing crash pad",
		Quantity:    2,
	},
	{
		Id:          1,
		Name:        "microspikes",
		Description: "",
		Quantity:    1,
	},
}

func (d *mockDB) GetItems() *[]api.Item {
	time.Sleep(time.Second * 1)

	clientData := mockItems

	return &clientData
}

func (d *mockDB) SetupDatabase() error {
	return nil
}
