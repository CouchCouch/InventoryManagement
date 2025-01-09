package tools

import (
	"database/sql"
	"fmt"
	"inventoryapi/api"
	"log"
	"time"

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

func (d *sqlDB) UpdateItem(item api.Item) bool {
    sql := "UPDATE items SET name=($1), description=($2), quantity=($3) WHERE id=($4)"

    rows, err := d.db.Query(sql, item.Name, item.Description, item.Quantity, item.Id)

    if err != nil {
        log.Fatalf("Failed to add item ERROR: %s", err)
    }

    rows.Close()

    return true
}

func (d *sqlDB) DeleteItem(id int) bool {
    sql := "DELETE FROM items WHERE id=(($1))"

    rows, err := d.db.Query(sql, id)

    if err != nil {
        log.Fatal("Failed to delete item")
        return false
    }

    rows.Close()

    return true
}

func (d *sqlDB) CheckoutItem(item api.CheckoutParams) *api.CheckoutItemReceipt {
    sql := "INSERT INTO checkouts (item_id, name, email, checkout_date) VALUES (($1), ($2), ($3), ($4))"

    checkoutDate := time.Now()

    rows, err := d.db.Query(sql, item.Id, item.Name, item.Email, checkoutDate)

    if err != nil {
        log.Fatal("Checkout Failed")
    }

    rows.Close()

    resp := api.CheckoutItemReceipt{
        ItemId: item.Id,
        Name: item.Name,
        Email: item.Email,
        Date: checkoutDate,
    }

    return &resp
}

func (d *sqlDB) ReturnItem(id int) bool {
    sql := "UPDATE CHECKOUTS SET returned = TRUE WHERE item_id=($1)"

    rows, err := d.db.Query(sql, id)

    if err != nil {
        log.Fatal("Return Failed")
    }

    rows.Close()

    return true;
}

func (d *sqlDB) GetCheckouts() *[]api.CheckoutItem {
    sql := "SELECT item_id, name, email, checkout_date, returned FROM checkouts"

    rows,  err := d.db.Query(sql)

    if err != nil {
        log.Fatal("Get Checkouts Failed", err)
    }

    defer rows.Close()

    var checkouts []api.CheckoutItem

    for rows.Next() {
        var itemId int
        var name string
        var email string
        var date time.Time
        var returned bool

        if err := rows.Scan(&itemId, &name, &email, &date, &returned); err != nil {
            log.Fatal(err)
        }

        checkouts = append(checkouts, api.CheckoutItem{
            ItemId: itemId,
            Name: name,
            Email: email,
            Date: date,
            Returned: returned,
        })
    }

    return &checkouts
}
