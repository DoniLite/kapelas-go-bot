package bot

import (
	"encoding/json"
	"slices"

	"github.com/DoniLite/kapelas-bot/core"
	"github.com/google/uuid"
)

type BotRepository struct {
	// Add any dependencies or fields needed for the service
}

func NewBotRepository() *BotRepository {
	return &BotRepository{
		// Initialize any dependencies or fields here
	}
}

// Add methods for your BotRepository here, e.g.:

func (r *BotRepository) GetProductList(categories ...string) ([]ShopProduct, error) {
	// Implement logic to retrieve product list from the store or an external API
	persistedProducts, err := core.GetStore().List(core.CollectionProducts.String())
	if err != nil {
		return []ShopProduct{}, err
	}
	products := make([]ShopProduct, len(persistedProducts))
	for _, m := range persistedProducts {
		b, _ := json.Marshal(m)
		var p ShopProduct
		if err := json.Unmarshal(b, &p); err != nil {
			continue
		}
		if len(categories) > 0 {
			included := slices.Contains(categories, p.CategoryID)
			if !included {
				continue
			}
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *BotRepository) CreateOrder(order ShopOrder) error {
	return core.GetStore().Create(core.CollectionOrders.String(), order.ID, order)
}

func (r *BotRepository) GetOrder(id string) (*ShopOrder, error) {
	var order ShopOrder
	err := core.GetStore().Get(core.CollectionOrders.String(), id, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *BotRepository) UpdateOrder(order ShopOrder) error {
	return core.GetStore().Upsert(core.CollectionOrders.String(), order.ID, order)
}

func (r *BotRepository) GetCategories() ([]ShopCategory, error) {
	var categories []ShopCategory
	persistedCategories, err := core.GetStore().List(core.CollectionProductCategories.String())
	if err != nil {
		return []ShopCategory{}, err
	}
	for _, m := range persistedCategories {
		b, _ := json.Marshal(m)
		var c ShopCategory
		if err := json.Unmarshal(b, &c); err != nil {
			continue
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *BotRepository) GetProductsByCategory(categoryID string) ([]ShopProduct, error) {
	var category ShopCategory
	err := core.GetStore().Get(core.CollectionProductCategories.String(), categoryID, &category)
	if err != nil {
		return []ShopProduct{}, err
	}
	return category.Products, nil
}

func (r *BotRepository) CreateCategory(category string) (string, error) {
	id := uuid.NewString()
	return id, core.GetStore().Create(core.CollectionProductCategories.String(), id, ShopCategory{
		ID:       id,
		Name:     category,
		Products: []ShopProduct{},
	})
}

func (r *BotRepository) AddProductToCategory(categoryID string, product ShopProduct) error {
	var category ShopCategory
	err := core.GetStore().Get(core.CollectionProductCategories.String(), categoryID, &category)
	if err != nil {
		return err
	}
	category.Products = append(category.Products, product)
	return core.GetStore().Upsert(core.CollectionProductCategories.String(), categoryID, category)
}

func (r *BotRepository) GetProductByID(productID string) (*ShopProduct, error) {
	var product ShopProduct
	err := core.GetStore().Get(core.CollectionProducts.String(), productID, &product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *BotRepository) GetProductCategory(categoryId string) (*ShopCategory, error) {
	var category ShopCategory
	err := core.GetStore().Get(core.CollectionProductCategories.String(), categoryId, &category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *BotRepository) CreateProduct(product ShopProduct) error {
	return core.GetStore().Create(core.CollectionProducts.String(), product.ID, product)
}

func (r *BotRepository) UpdateProduct(product ShopProduct) error {
	return core.GetStore().Upsert(core.CollectionProducts.String(), product.ID, product)
}
