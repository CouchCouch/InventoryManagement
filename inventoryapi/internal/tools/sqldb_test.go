package tools

import (
	"context"
	"inventoryapi/api"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testDb *sqlDB

func TestMain(m *testing.M) {
	ctx := context.Background()
	connStr, cleanup := CreateTestDatabase(ctx)
	defer cleanup()

	var err error
	testDb, err = New(connStr)
	if err != nil {
		log.Fatalf("failed to create test database: %s", err)
	}

	os.Exit(m.Run())
}

func Test_sqlDB_GetItems(t *testing.T) {
	tests := []struct {
		name    string
		want    *[]api.Item
		wantErr bool
	}{
		{
			name: "get all items from a clean database",
			want: &[]api.Item{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testDb.GetItems()
			if (err != nil) != tt.wantErr {
				t.Errorf("sqlDB.GetItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(*got) != len(*tt.want) {
				t.Errorf("sqlDB.GetItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sqlDB_AddItem_and_GetItem(t *testing.T) {
	newItem := api.NewItem{
		Name:        "Test Item",
		Description: "A test item",
		Quantity:    1,
	}

	// Add the item
	id, err := testDb.AddItem(newItem)
	if err != nil {
		t.Fatalf("AddItem() error = %v", err)
	}

	// Get the item
	items, err := testDb.GetItem(*id)
	if err != nil {
		t.Fatalf("GetItem() error = %v", err)
	}

	// Check that we got one item back
	if len(*items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(*items))
	}

	// Check that the item is correct
	item := (*items)[0]
	if item.Name != newItem.Name || item.Description != newItem.Description || item.Quantity != newItem.Quantity {
		t.Errorf("Got item %v, want %v", item, newItem)
	}

	// Clean up the item
	err = testDb.DeleteItem(*id)
	if err != nil {
		t.Fatalf("DeleteItem() error = %v", err)
	}
}
