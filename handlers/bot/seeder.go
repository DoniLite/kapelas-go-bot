package bot

import "github.com/google/uuid"

var products map[string][]struct {
	name        string
	description string
	price       float64
	images      []string
} = map[string][]struct {
	name        string
	description string
	price       float64
	images      []string
}{
	"electronics": {
		{
			name:        "Smartphone",
			description: "A high-end smartphone with a sleek design and powerful features.",
			price:       699.99,
			images: []string{
				"https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?auto=format&fit=crop&w=1200&q=80",
				"https://images.unsplash.com/photo-1598327105666-5b89351aff97?auto=format&fit=crop&w=1200&q=80",
			},
		},
		{
			name:        "Laptop",
			description: "A lightweight laptop with a long battery life, perfect for work and entertainment.",
			price:       1299.99,
			images: []string{
				"https://images.unsplash.com/photo-1496181133206-80ce9b88a853?auto=format&fit=crop&w=1200&q=80",
				"https://images.unsplash.com/photo-1517336714731-489689fd1ca8?auto=format&fit=crop&w=1200&q=80",
			},
		},
		{
			name:        "Tablet",
			description: "A versatile tablet with a large display and powerful performance.",
			price:       499.99,
			images: []string{
				"https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?auto=format&fit=crop&w=1200&q=80",
				"https://images.unsplash.com/photo-1561154464-82e9adf32764?auto=format&fit=crop&w=1200&q=80",
			},
		},
	},
	"books": {
		{
			name:        "The Great Gatsby",
			description: "A classic novel by F. Scott Fitzgerald, set in the Roaring Twenties.",
			price:       10.99,
			images: []string{
				"https://images.unsplash.com/photo-1544947950-fa07a98d237f?auto=format&fit=crop&w=1200&q=80",
				"https://images.unsplash.com/photo-1512820790803-83ca734da794?auto=format&fit=crop&w=1200&q=80",
			},
		},
		{
			name:        "To Kill a Mockingbird",
			description: "A powerful novel by Harper Lee that explores themes of racial injustice and moral growth.",
			price:       8.99,
			images: []string{
				"https://images.unsplash.com/photo-1524995997946-a1c2e315a42f?auto=format&fit=crop&w=1200&q=80",
				"https://images.unsplash.com/photo-1456513080510-7bf3a84b82f8?auto=format&fit=crop&w=1200&q=80",
			},
		},
		{
			name:        "1984",
			description: "A dystopian novel by George Orwell that delves into themes of totalitarianism and surveillance.",
			price:       9.99,
			images: []string{
				"https://images.unsplash.com/photo-1516979187457-637abb4f9353?auto=format&fit=crop&w=1200&q=80",
				"https://images.unsplash.com/photo-1495446815901-a7297e633e8d?auto=format&fit=crop&w=1200&q=80",
			},
		},
	},
	"clothing": {
		{
			name:        "Cotton T-Shirt",
			description: "A comfortable and stylish cotton t-shirt for everyday wear.",
			price:       19.99,
			images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?auto=format&fit=crop&w=1200&q=80",
				"https://images.unsplash.com/photo-1503341504253-dff4815485f1?auto=format&fit=crop&w=1200&q=80",
			},
		},
		{
			name:        "Jeans",
			description: "Classic denim jeans with a modern fit, perfect for any occasion.",
			price:       49.99,
			images: []string{
				"https://images.unsplash.com/photo-1542272604-787c3835535d?auto=format&fit=crop&w=1200&q=80",
				"https://images.unsplash.com/photo-1511105043137-7e66f28270e3?auto=format&fit=crop&w=1200&q=80",
			},
		},
		{
			name:        "Jacket",
			description: "A stylish and warm jacket for the colder months.",
			price:       89.99,
			images: []string{
				"https://images.unsplash.com/photo-1551028719-00167b16eac5?auto=format&fit=crop&w=1200&q=80",
				"https://images.unsplash.com/photo-1520975954732-35dd22299614?auto=format&fit=crop&w=1200&q=80",
			},
		},
	},
	"home": {
		{
			name:        "Ceramic Vase",
			description: "A beautifully crafted ceramic vase to enhance your home decor.",
			price:       49.99,
			images: []string{
				"https://images.unsplash.com/photo-1612196808214-b8e1d6145a8c?auto=format&fit=crop&w=1200&q=80",
				"https://images.unsplash.com/photo-1616486338812-3dadae4b4ace?auto=format&fit=crop&w=1200&q=80",
			},
		},
		{
			name:        "Wall Art",
			description: "A stunning piece of wall art to add a touch of elegance to any room.",
			price:       79.99,
			images: []string{
				"https://images.unsplash.com/photo-1547891654-e66ed7ebb968?auto=format&fit=crop&w=1200&q=80",
				"https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=1200&q=80",
			},
		},
		{
			name:        "Throw Pillow",
			description: "A cozy throw pillow to add comfort and style to your living space.",
			price:       29.99,
			images: []string{
				"https://images.unsplash.com/photo-1584100936595-c0654b55a2e6?auto=format&fit=crop&w=1200&q=80",
				"https://images.unsplash.com/photo-1618220179428-22790b461013?auto=format&fit=crop&w=1200&q=80",
			},
		},
	},
	"toys": {
		{
			name:        "Building Blocks Set",
			description: "A fun and educational building blocks set for children of all ages.",
			price:       29.99,
			images: []string{
				"https://images.unsplash.com/photo-1587654780291-39c9404d746b?auto=format&fit=crop&w=1200&q=80",
				"https://images.unsplash.com/photo-1558060370-d644479cb6f7?auto=format&fit=crop&w=1200&q=80",
			},
		},
		{
			name:        "Remote Control Car",
			description: "A fast and agile remote control car for hours of entertainment.",
			price:       49.99,
			images: []string{
				"https://images.unsplash.com/photo-1594787318286-3d835c1d207f?auto=format&fit=crop&w=1200&q=80",
				"https://images.unsplash.com/photo-1503736334956-4c8f8e92946d?auto=format&fit=crop&w=1200&q=80",
			},
		},
		{
			name:        "Dollhouse",
			description: "A charming dollhouse with intricate details and accessories.",
			price:       89.99,
			images: []string{
				"https://images.unsplash.com/photo-1604881988758-f76ad2f7aac1?auto=format&fit=crop&w=1200&q=80",
				"https://images.unsplash.com/photo-1596461404969-9ae70f2830c1?auto=format&fit=crop&w=1200&q=80",
			},
		},
	},
}

func SeedProducts() {
	repo := NewBotRepository()
	for category, details := range products {
		categoryID, err := repo.CreateCategory(category)
		if err != nil {
			continue
		}
		for _, product := range details {
			repo.CreateProduct(ShopProduct{
				ID:          uuid.NewString(),
				Name:        product.name,
				Description: product.description,
				Price:       product.price,
				Images:      product.images,
				CategoryID:  categoryID,
			})
		}
	}
}
