package tools

import (
    "time"
)

type mockDB struct{}

var mockLoginDetails = map[string]LoginDetails{
    "ryan": {
        AuthToken: "abc123",
        Username: "ryan",
    },
    "jackson": {
        AuthToken: "123abc",
        Username: "jackson",
    },
}

var mockItems = map[string]UserItemDetails{
    "ryan": {
        Id: 2,
        Name: "crach pad",
        Description: "climbing crash pad",
        Quantity: 2,
        CheckoutDate: time.Date(2024, time.December, 17, 22, 0, 0, 0, time.UTC),
    },
    "jackson": {
        Id: 1,
        Name: "microspikes",
        Description: "",
        Quantity: 1,
        CheckoutDate: time.Date(2024, time.December, 17, 22, 0, 0, 0, time.UTC),
    },
}

func (d *mockDB) GetUserLoginDetails(username string) *LoginDetails {
    time.Sleep(time.Second * 1)

    var clientData = LoginDetails{}
    clientData, ok := mockLoginDetails[username]
    if !ok {
        return nil
    }

    return &clientData
}

func (d *mockDB) GetUserItems(username string) *UserItemDetails {
    time.Sleep(time.Second * 1)

    var clientData = UserItemDetails{}
    clientData, ok := mockItems[username]
    if !ok {
        return nil
    }

    return &clientData
}

func (d *mockDB) SetupDatabase() error {
    return nil
}
