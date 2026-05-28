package bot

import (
	"fmt"
	"math"
	"strings"

	"github.com/Arnel7/kappelas-sdk-go"
	"github.com/DoniLite/kapelas-bot/conf"
	"github.com/google/uuid"
)

type BotService struct {
	bot        *kappelas.Bot
	repository *BotRepository
}

func NewBotService(bot *kappelas.Bot, repository *BotRepository) *BotService {
	return &BotService{
		bot:        bot,
		repository: repository,
	}
}

func (s *BotService) ListAllAvailableCategories() (*kappelas.InlineKeyboard, error) {
	categories, err := s.repository.GetCategories()
	if err != nil {
		return nil, err
	}
	return GenerateProductCategoryListMarker(categories...), nil
}

func (s *BotService) ListAvailableProducts(categoryId string, page int) (*kappelas.InlineKeyboard, ProductListPage, error) {
	if page < 1 {
		page = 1
	}

	products, err := s.repository.GetProductList(categoryId)
	if err != nil {
		return nil, ProductListPage{}, err
	}

	total := len(products)
	totalPages := int(math.Ceil(float64(total) / float64(ProductListPageSize)))
	if totalPages == 0 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * ProductListPageSize
	end := start + ProductListPageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	pageData := ProductListPage{
		Products:   products[start:end],
		CategoryID: categoryId,
		Page:       page,
		PageSize:   ProductListPageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	return GenerateProductListMarker(pageData), pageData, nil
}

func (s *BotService) GetCategoryById(categoryId string) (*ShopCategory, error) {
	return s.repository.GetProductCategory(categoryId)
}

func (s *BotService) GetProductById(productId string) (*ShopProduct, error) {
	return s.repository.GetProductByID(productId)
}

const (
	stateOwnerAccessToken = "owner_access_token"
	stateOwnerChat        = "owner_chat"
	stateMessagePrefix    = "message:"
)

func (s *BotService) PlaceOrder(userID string, chatID int64, productID string, quantity int) (*ShopOrder, *ShopProduct, error) {
	if quantity < 1 {
		quantity = 1
	}
	product, err := s.repository.GetProductByID(productID)
	if err != nil {
		return nil, nil, err
	}
	order := ShopOrder{
		ID:         uuid.NewString(),
		UserID:     userID,
		ChatID:     chatID,
		ProductID:  product.ID,
		Quantity:   quantity,
		TotalPrice: TotalizePrice(product.Price, quantity),
		Status:     OrderStatusPending,
	}
	if err := s.repository.CreateOrder(order); err != nil {
		return nil, nil, err
	}
	return &order, product, nil
}

func (s *BotService) ListOrdersByUser(userID string) ([]ShopOrder, error) {
	return s.repository.GetOrdersByUser(userID)
}

func (s *BotService) ListOrders() ([]ShopOrder, error) {
	return s.repository.GetOrders()
}

func (s *BotService) UpdateOrderStatus(orderID string, status OrderStatus) (*ShopOrder, error) {
	if !status.IsValid() {
		return nil, fmt.Errorf("invalid order status: %s", status)
	}
	order, err := s.repository.GetOrder(orderID)
	if err != nil {
		return nil, err
	}
	order.Status = status
	if err := s.repository.UpdateOrder(*order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *BotService) ListCategories() ([]ShopCategory, error) {
	return s.repository.GetCategories()
}

func (s *BotService) CreateCategory(name string) (*ShopCategory, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("category name is required")
	}
	id, err := s.repository.CreateCategory(name)
	if err != nil {
		return nil, err
	}
	return s.repository.GetProductCategory(id)
}

func (s *BotService) UpdateCategory(categoryID string, name string) (*ShopCategory, error) {
	category, err := s.repository.GetProductCategory(categoryID)
	if err != nil {
		return nil, err
	}
	category.Name = strings.TrimSpace(name)
	if category.Name == "" {
		return nil, fmt.Errorf("category name is required")
	}
	if err := s.repository.UpdateCategory(*category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *BotService) DeleteCategory(categoryID string) error {
	if _, err := s.repository.GetProductCategory(categoryID); err != nil {
		return err
	}
	products, err := s.repository.GetProductList(categoryID)
	if err != nil {
		return err
	}
	for _, product := range products {
		if err := s.repository.DeleteProduct(product.ID); err != nil {
			return err
		}
	}
	return s.repository.DeleteCategory(categoryID)
}

func (s *BotService) CreateProduct(product ShopProduct) (*ShopProduct, error) {
	if strings.TrimSpace(product.Name) == "" {
		return nil, fmt.Errorf("product name is required")
	}
	if strings.TrimSpace(product.CategoryID) == "" {
		return nil, fmt.Errorf("product category is required")
	}
	if product.Price < 0 {
		return nil, fmt.Errorf("product price cannot be negative")
	}
	if _, err := s.repository.GetProductCategory(product.CategoryID); err != nil {
		return nil, err
	}
	product.ID = uuid.NewString()
	if err := s.repository.CreateProduct(product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *BotService) UpdateProduct(productID string, updates map[string]string) (*ShopProduct, error) {
	product, err := s.repository.GetProductByID(productID)
	if err != nil {
		return nil, err
	}
	if name := strings.TrimSpace(updates["name"]); name != "" {
		product.Name = name
	}
	if description := strings.TrimSpace(updates["description"]); description != "" {
		product.Description = description
	}
	if categoryID := strings.TrimSpace(updates["category"]); categoryID != "" {
		if _, err := s.repository.GetProductCategory(categoryID); err != nil {
			return nil, err
		}
		product.CategoryID = categoryID
	}
	if images := strings.TrimSpace(updates["images"]); images != "" {
		product.Images = splitCSV(images)
	}
	if price := strings.TrimSpace(updates["price"]); price != "" {
		var value float64
		if _, err := fmt.Sscanf(price, "%f", &value); err != nil {
			return nil, fmt.Errorf("invalid product price")
		}
		if value < 0 {
			return nil, fmt.Errorf("product price cannot be negative")
		}
		product.Price = value
	}
	if err := s.repository.UpdateProduct(*product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *BotService) DeleteProduct(productID string) error {
	return s.repository.DeleteProduct(productID)
}

func (s *BotService) EnsureOwnerAccessToken() (string, error) {
	if token := strings.TrimSpace(conf.GetEnv().GetString(conf.OWNER_ACCESS_TOKEN)); token != "" {
		return token, nil
	}
	item, err := s.repository.GetState(stateOwnerAccessToken)
	if err == nil && strings.TrimSpace(item.Value) != "" {
		return item.Value, nil
	}
	token := uuid.NewString()
	err = s.repository.SetState(BotStateItem{ID: stateOwnerAccessToken, Value: token})
	return token, err
}

func (s *BotService) SaveOwnerChat(chatID int64) error {
	return s.repository.SetState(BotStateItem{ID: stateOwnerChat, ChatID: chatID})
}

func (s *BotService) GetOwnerChat() (int64, bool) {
	item, err := s.repository.GetState(stateOwnerChat)
	if err != nil || item.ChatID == 0 {
		return 0, false
	}
	return item.ChatID, true
}

func (s *BotService) IsOwnerChat(chatID int64) bool {
	ownerChatID, ok := s.GetOwnerChat()
	return ok && ownerChatID == chatID
}

func (s *BotService) SaveMessageRef(chatID int64, key string, messageID int64) error {
	id := stateMessagePrefix + fmt.Sprintf("%d:%s", chatID, key)
	return s.repository.SetState(BotStateItem{ID: id, ChatID: chatID, Message: messageID})
}

func (s *BotService) GetMessageRef(chatID int64, key string) (int64, bool) {
	id := stateMessagePrefix + fmt.Sprintf("%d:%s", chatID, key)
	item, err := s.repository.GetState(id)
	if err != nil || item.Message == 0 {
		return 0, false
	}
	return item.Message, true
}

func (status OrderStatus) IsValid() bool {
	switch status {
	case OrderStatusPending, OrderStatusConfirmed, OrderStatusShipped, OrderStatusDelivered, OrderStatusCancelled:
		return true
	default:
		return false
	}
}

func splitCSV(value string) []string {
	var result []string
	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
