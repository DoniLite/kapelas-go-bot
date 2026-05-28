package bot

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type ShopOrder struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id"`
	ChatID     int64       `json:"chat_id"`
	ProductID  string      `json:"product_id"`
	Quantity   int         `json:"quantity,omitempty"`
	TotalPrice float64     `json:"total_price"`
	Status     OrderStatus `json:"status"`
}

type ShopProduct struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Images      []string `json:"images"`
	CategoryID  string   `json:"category_id"`
}

type ShopCategory struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Products []ShopProduct `json:"products"`
}

type BotStateItem struct {
	ID      string `json:"id"`
	Value   string `json:"value"`
	ChatID  int64  `json:"chat_id,omitempty"`
	Message int64  `json:"message,omitempty"`
}
