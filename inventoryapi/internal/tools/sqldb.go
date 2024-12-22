package tools

import (
	"database/sql"
	"inventoryapi/api"
	"log"

	_ "github.com/lib/pq"
)

func getItems(db *sql.DB) []api.Item {
    sql := "SELECT id, name, description, quantity FROM items"

    rows, err := db.Query(sql)

    if err != nil {
        log.Fatalf("Failed to fetch items")
    }

    defer rows.Close()

    var items []api.Item

    for rows.Next() {
        var itemId int
        var itemName string
        var itemDescription string
        var itemQuantity int

        if err := rows.Scan(&itemId, &itemName, &itemDescription, &itemQuantity); err != nil {
            log.Fatal(err)
        }

        items = append(items, api.Item{
            Id: itemId,
            Name: itemName,
            Description: itemDescription,
            Quantity: itemQuantity,
        })
    }

    return items
}

func getItem(id int, db *sql.DB) api.Item{
    sql := "SELECT id, name, description, quantity FROM items where id=($1)"

    rows, err := db.Query(sql, id)

    if err != nil {
        log.Fatalf("Failed to fetch items")
    }

    defer rows.Close()


    rows.Next()
    var itemId int
    var itemName string
    var itemDescription string
    var itemQuantity int

    if err := rows.Scan(&itemId, &itemName, &itemDescription, &itemQuantity); err != nil {
        log.Fatal(err)
    }

    item := api.Item{
        Id: itemId,
        Name: itemName,
        Description: itemDescription,
        Quantity: itemQuantity,
    }

    return item
}

