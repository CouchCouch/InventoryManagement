package tools

import (
	"database/sql"
	"fmt"
	"inventoryapi/api"
	"log"

	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

type sqlDB struct{
    db *sql.DB
}

func (d *sqlDB) GetItems() *[]api.Item {
    sql := "SELECT id, name, description, quantity FROM items"

    rows, err := d.db.Query(sql)

    if err != nil {
        log.Fatal(err)
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

    return &items
}

type sqlDBCredentials struct {
    Host string `yaml:"host"`
    Database string `yaml:"database"`
    User string `yaml:"user"`
    Password string `yaml:"password"`
    Port int `yaml:"port"`
}

func (d *sqlDB) SetupDatabase() error {
    var credentials sqlDBCredentials

    yamlFile, err := os.ReadFile("../.config/db.yml")
    if err != nil {
        panic(err)
    }

    err = yaml.Unmarshal(yamlFile, &credentials)
    if err != nil {
        panic(err)
    }

    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", credentials.Host, credentials.Port, credentials.User, credentials.Password, credentials.Database)

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    d.db = db

    return nil
}

func (d *sqlDB) CloseDatabase() error {
    return d.db.Close()
}

func (d *sqlDB) GetItem(id int) *[]api.Item{
    sql := "SELECT id, name, description, quantity FROM items where id=($1)"

    rows, err := d.db.Query(sql, id)

    if err != nil {
        log.Fatal("Failed to fetch item")
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

    item := []api.Item{{
        Id: itemId,
        Name: itemName,
        Description: itemDescription,
        Quantity: itemQuantity,
    }}

    return &item
}

func (d *sqlDB) AddItem(item api.NewItem) *int {
    sql := "INSERT INTO items (name, description, quantity) VALUES (($1), ($2), ($3)) RETURNING id"

    rows, err := d.db.Query(sql, item.Name, item.Description, item.Quantity)

    if err != nil {
        log.Fatal("Failed to add item")
    }

    defer rows.Close()

    rows.Next()
    var itemId int

    if err := rows.Scan(&itemId); err != nil {
        log.Fatal(err)
    }

    return &itemId
}
