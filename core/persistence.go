package core

import (
	"log"

	"github.com/DoniLite/kapelas-bot/core/store"
)

var (
	// Store is the global instance of the JSONStore used for data persistence.
	global_store *store.JSONStore
)

type Collection string

func (c Collection) String() string {
	return string(c)
}

const (
	// Collection names for the JSON store.
	CollectionProducts          Collection = "products"
	CollectionOrders            Collection = "orders"
	CollectionProductCategories Collection = "product_categories"
)

func GetStore() *store.JSONStore {
	return global_store
}

func init() {
	var err error
	global_store, err = store.NewJSONStore("")
	if err != nil {
		log.Fatalf("Failed to initialize JSON store: %s", err.Error())
	}
}
