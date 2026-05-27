package bot

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
			images:      []string{"https://example.com/images/smartphone1.jpg", "https://example.com/images/smartphone2.jpg"},
		},
		{
			name:        "Laptop",
			description: "A lightweight laptop with a long battery life, perfect for work and entertainment.",
			price:       1299.99,
			images:      []string{"https://example.com/images/laptop1.jpg", "https://example.com/images/laptop2.jpg"},
		},
		{
			name:        "Tablet",
			description: "A versatile tablet with a large display and powerful performance.",
			price:       499.99,
			images:      []string{"https://example.com/images/tablet1.jpg", "https://example.com/images/tablet2.jpg"},
		},
	},
	"books": {
		{
			name:        "The Great Gatsby",
			description: "A classic novel by F. Scott Fitzgerald, set in the Roaring Twenties.",
			price:       10.99,
			images:      []string{"https://example.com/images/gatsby1.jpg", "https://example.com/images/gatsby2.jpg"},
		},
		{
			name:        "To Kill a Mockingbird",
			description: "A powerful novel by Harper Lee that explores themes of racial injustice and moral growth.",
			price:       8.99,
			images:      []string{"https://example.com/images/mockingbird1.jpg", "https://example.com/images/mockingbird2.jpg"},
		},
		{
			name:        "1984",
			description: "A dystopian novel by George Orwell that delves into themes of totalitarianism and surveillance.",
			price:       9.99,
			images:      []string{"https://example.com/images/1984_1.jpg", "https://example.com/images/1984_2.jpg"},
		},
	},
	"clothing": {
		{
			name:        "Cotton T-Shirt",
			description: "A comfortable and stylish cotton t-shirt for everyday wear.",
			price:       19.99,
			images:      []string{"https://example.com/images/tshirt1.jpg", "https://example.com/images/tshirt2.jpg"},
		},
		{
			name:        "Jeans",
			description: "Classic denim jeans with a modern fit, perfect for any occasion.",
			price:       49.99,
			images:      []string{"https://example.com/images/jeans1.jpg", "https://example.com/images/jeans2.jpg"},
		},
		{
			name:        "Jacket",
			description: "A stylish and warm jacket for the colder months.",
			price:       89.99,
			images:      []string{"https://example.com/images/jacket1.jpg", "https://example.com/images/jacket2.jpg"},
		},
	},
	"home": {
		{
			name:        "Ceramic Vase",
			description: "A beautifully crafted ceramic vase to enhance your home decor.",
			price:       49.99,
			images:      []string{"https://example.com/images/vase1.jpg", "https://example.com/images/vase2.jpg"},
		},
		{
			name:        "Wall Art",
			description: "A stunning piece of wall art to add a touch of elegance to any room.",
			price:       79.99,
			images:      []string{"https://example.com/images/wallart1.jpg", "https://example.com/images/wallart2.jpg"},
		},
		{
			name:        "Throw Pillow",
			description: "A cozy throw pillow to add comfort and style to your living space.",
			price:       29.99,
			images:      []string{"https://example.com/images/pillow1.jpg", "https://example.com/images/pillow2.jpg"},
		},
	},
	"toys": {
		{
			name:        "Building Blocks Set",
			description: "A fun and educational building blocks set for children of all ages.",
			price:       29.99,
			images:      []string{"https://example.com/images/blocks1.jpg", "https://example.com/images/blocks2.jpg"},
		},
		{
			name:        "Remote Control Car",
			description: "A fast and agile remote control car for hours of entertainment.",
			price:       49.99,
			images:      []string{"https://example.com/images/rccar1.jpg", "https://example.com/images/rccar2.jpg"},
		},
		{
			name:        "Dollhouse",
			description: "A charming dollhouse with intricate details and accessories.",
			price:       89.99,
			images:      []string{"https://example.com/images/dollhouse1.jpg", "https://example.com/images/dollhouse2.jpg"},
		},
	},
}

func SeedProducts() {
	repo := NewBotRepository()
	for category, details := range products {
		productId, err := repo.CreateCategory(category)
		if err != nil {
			continue
		}
		for _, product := range details {
			repo.CreateProduct(ShopProduct{
				ID:          productId,
				Name:        product.name,
				Description: product.description,
				Price:       product.price,
				Images:      product.images,
				CategoryID:  productId,
			})
		}
	}
}
