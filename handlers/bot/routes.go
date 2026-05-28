package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/Arnel7/kappelas-sdk-go"
	"github.com/DoniLite/kapelas-bot/core"
)

func HandleStartCommand(bot *kappelas.Bot, service *BotService, chatId int64, senderName string) {
	ctx := context.Background()
	markup, err := service.ListAllAvailableCategories()
	if err != nil {
		bot.Messages.Send(ctx, kappelas.SendMessageParams{
			ChatID: chatId,
			Text:   core.STATIC_FALLBACK_TEXT,
		})
		return
	}
	sendOrEdit(bot, service, chatId, "start", fmt.Sprintf(core.STATIC_WELCOME_TEXT, senderName, "doni"), markup)
}

func HandleHelpCommand(bot *kappelas.Bot, chatId int64) {
	ctx := context.Background()
	bot.Messages.Send(ctx, kappelas.SendMessageParams{
		ChatID:         chatId,
		Text:           core.STATIC_HELP_TEXT,
		DeletePrevious: true,
	})
}

func HandleListProducts(bot *kappelas.Bot, service *BotService, chatId int64, query map[string]string) {
	ctx := context.Background()
	categoryId := query["category"]
	if categoryId == "" {
		markup, err := service.ListAllAvailableCategories()
		if err != nil {
			bot.Messages.Send(ctx, kappelas.SendMessageParams{
				ChatID: chatId,
				Text:   core.STATIC_FALLBACK_TEXT,
			})
			return
		}
		sendOrEdit(bot, service, chatId, "categories", core.STATIC_CHOOSE_CATEGORY_TEXT, markup)
		return
	}

	page, err := strconv.Atoi(query["page"])
	if err != nil || page < 1 {
		page = 1
	}
	log.Printf("Listing products for category id: %s", categoryId)
	category, err := service.GetCategoryById(categoryId)
	if err != nil {
		bot.Messages.Send(ctx, kappelas.SendMessageParams{
			ChatID: chatId,
			Text:   core.STATIC_FALLBACK_TEXT,
		})
		return
	}
	markup, pageData, err := service.ListAvailableProducts(categoryId, page)
	if err != nil {
		bot.Messages.Send(ctx, kappelas.SendMessageParams{
			ChatID: chatId,
			Text:   core.STATIC_FALLBACK_TEXT,
		})
		return
	}
	text := fmt.Sprintf(core.STATIC_LIST_PRODUCTS_BY_CATEGORY_PAGE_TEXT, category.Name, pageData.Page, pageData.TotalPages)
	if pageData.Total == 0 {
		text = fmt.Sprintf(core.STATIC_EMPTY_PRODUCT_CATEGORY_TEXT, category.Name)
	}
	sendOrEdit(bot, service, chatId, "products:"+categoryId, text, markup)
}

func HandleViewProduct(bot *kappelas.Bot, service *BotService, chatId int64, query map[string]string) {
	ctx := context.Background()
	productId := query["product"]
	if productId == "" {
		bot.Messages.Send(ctx, kappelas.SendMessageParams{
			ChatID: chatId,
			Text:   core.STATIC_FALLBACK_TEXT,
		})
		return
	}

	product, err := service.GetProductById(productId)
	if err != nil {
		bot.Messages.Send(ctx, kappelas.SendMessageParams{
			ChatID: chatId,
			Text:   core.STATIC_FALLBACK_TEXT,
		})
		return
	}

	orderCallback := BuildCallbackData(PlaceOrderCommand, map[string]string{
		"product": product.ID,
	})
	markup := &kappelas.InlineKeyboard{
		InlineKeyboard: [][]kappelas.InlineKeyboardButton{
			{
				{
					Text:         "Place order",
					CallbackData: &orderCallback,
				},
			},
		},
	}

	text := fmt.Sprintf(core.STATIC_VIEW_PRODUCT_TEXT, product.Name, product.Description, product.Price, product.ID, product.CategoryID)
	sendOrEdit(bot, service, chatId, "product:"+product.ID, text, markup)
	sendProductImages(bot, chatId, product)
}

func HandlePlaceOrder(bot *kappelas.Bot, service *BotService, chatId int64, userID string, query map[string]string) {
	ctx := context.Background()
	productID := strings.TrimSpace(query["product"])
	if productID == "" {
		sendBadUsage(bot, chatId)
		return
	}
	quantity, err := strconv.Atoi(query["quantity"])
	if err != nil || quantity < 1 {
		quantity = 1
	}

	order, product, err := service.PlaceOrder(userID, chatId, productID, quantity)
	if err != nil {
		bot.Messages.Send(ctx, kappelas.SendMessageParams{
			ChatID: chatId,
			Text:   core.STATIC_FALLBACK_TEXT,
		})
		return
	}
	bot.Messages.Send(ctx, kappelas.SendMessageParams{
		ChatID:         chatId,
		DeletePrevious: true,
		Text: fmt.Sprintf(
			core.STATIC_ORDER_CREATED_TEXT,
			order.ID,
			product.Name,
			order.Quantity,
			order.TotalPrice,
			order.Status,
		),
	})
	notifyOwner(bot, service, fmt.Sprintf("new order %s for product %s by user %s", order.ID, product.ID, userID))
}

func HandleMyOrders(bot *kappelas.Bot, service *BotService, chatId int64, userID string) {
	ctx := context.Background()
	orders, err := service.ListOrdersByUser(userID)
	if err != nil {
		bot.Messages.Send(ctx, kappelas.SendMessageParams{
			ChatID: chatId,
			Text:   core.STATIC_FALLBACK_TEXT,
		})
		return
	}
	sendOrEdit(bot, service, chatId, "my-orders", formatOrders(orders), nil)
}

func HandleRequestOwnerAccess(bot *kappelas.Bot, service *BotService, chatId int64, token string) {
	ctx := context.Background()
	expectedToken, err := service.EnsureOwnerAccessToken()
	if err != nil {
		bot.Messages.Send(ctx, kappelas.SendMessageParams{
			ChatID: chatId,
			Text:   core.STATIC_FALLBACK_TEXT,
		})
		return
	}
	if expectedToken == "" || strings.TrimSpace(token) != expectedToken {
		bot.Messages.Send(ctx, kappelas.SendMessageParams{
			ChatID: chatId,
			Text:   core.STATIC_UNAUTHORIZED_TEXT,
		})
		return
	}
	if err := service.SaveOwnerChat(chatId); err != nil {
		bot.Messages.Send(ctx, kappelas.SendMessageParams{
			ChatID: chatId,
			Text:   core.STATIC_FALLBACK_TEXT,
		})
		return
	}
	bot.Messages.Send(ctx, kappelas.SendMessageParams{
		ChatID: chatId,
		Text:   fmt.Sprintf(core.STATIC_OWNER_ACCESS_GRANTED_TEXT, chatId),
	})
}

func HandleAdminViewCategories(bot *kappelas.Bot, service *BotService, chatId int64, isAdmin bool) {
	if !isAdmin {
		sendUnauthorized(bot, chatId)
		return
	}
	ctx := context.Background()
	categories, err := service.ListCategories()
	if err != nil {
		bot.Messages.Send(ctx, kappelas.SendMessageParams{
			ChatID: chatId,
			Text:   core.STATIC_FALLBACK_TEXT,
		})
		return
	}
	sendOrEdit(bot, service, chatId, "admin-categories", formatCategories(categories), nil)
}

func HandleOwnerAccessToken(bot *kappelas.Bot, service *BotService, chatId int64, isAdmin bool) {
	if !isAdmin {
		sendUnauthorized(bot, chatId)
		return
	}
	token, err := service.EnsureOwnerAccessToken()
	if err != nil {
		sendText(bot, chatId, core.STATIC_FALLBACK_TEXT)
		return
	}
	sendText(bot, chatId, fmt.Sprintf(core.STATIC_OWNER_ACCESS_TOKEN_TEXT, token))
}

func HandleAddCategory(bot *kappelas.Bot, service *BotService, chatId int64, isAdmin bool, args string) {
	if !isAdmin {
		sendUnauthorized(bot, chatId)
		return
	}
	category, err := service.CreateCategory(args)
	if err != nil {
		sendBadUsage(bot, chatId)
		return
	}
	sendText(bot, chatId, fmt.Sprintf(core.STATIC_CATEGORY_CREATED_TEXT, category.ID, category.Name))
	notifyOwner(bot, service, fmt.Sprintf("category created: %s (%s)", category.Name, category.ID))
}

func HandleUpdateCategory(bot *kappelas.Bot, service *BotService, chatId int64, isAdmin bool, args string) {
	if !isAdmin {
		sendUnauthorized(bot, chatId)
		return
	}
	parts := splitPipeArgs(args)
	if len(parts) != 2 {
		sendBadUsage(bot, chatId)
		return
	}
	category, err := service.UpdateCategory(parts[0], parts[1])
	if err != nil {
		sendBadUsage(bot, chatId)
		return
	}
	sendText(bot, chatId, fmt.Sprintf(core.STATIC_CATEGORY_UPDATED_TEXT, category.ID, category.Name))
	notifyOwner(bot, service, fmt.Sprintf("category updated: %s (%s)", category.Name, category.ID))
}

func HandleDeleteCategory(bot *kappelas.Bot, service *BotService, chatId int64, isAdmin bool, args string) {
	if !isAdmin {
		sendUnauthorized(bot, chatId)
		return
	}
	categoryID := strings.TrimSpace(args)
	if categoryID == "" {
		sendBadUsage(bot, chatId)
		return
	}
	if err := service.DeleteCategory(categoryID); err != nil {
		sendBadUsage(bot, chatId)
		return
	}
	sendText(bot, chatId, fmt.Sprintf(core.STATIC_CATEGORY_DELETED_TEXT, categoryID))
	notifyOwner(bot, service, fmt.Sprintf("category deleted: %s", categoryID))
}

func HandleAddProduct(bot *kappelas.Bot, service *BotService, chatId int64, isAdmin bool, args string) {
	if !isAdmin {
		sendUnauthorized(bot, chatId)
		return
	}
	parts := splitPipeArgs(args)
	if len(parts) < 4 {
		sendBadUsage(bot, chatId)
		return
	}
	price, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		sendBadUsage(bot, chatId)
		return
	}
	images := []string{}
	if len(parts) > 4 {
		images = splitCSV(parts[4])
	}
	product, err := service.CreateProduct(ShopProduct{
		CategoryID:  parts[0],
		Name:        parts[1],
		Price:       price,
		Description: parts[3],
		Images:      images,
	})
	if err != nil {
		sendBadUsage(bot, chatId)
		return
	}
	sendText(bot, chatId, fmt.Sprintf(core.STATIC_PRODUCT_CREATED_TEXT, product.ID, product.Name, product.Price))
	notifyOwner(bot, service, fmt.Sprintf("product created: %s (%s)", product.Name, product.ID))
}

func HandleUpdateProduct(bot *kappelas.Bot, service *BotService, chatId int64, isAdmin bool, args string) {
	if !isAdmin {
		sendUnauthorized(bot, chatId)
		return
	}
	parts := splitPipeArgs(args)
	if len(parts) < 2 {
		sendBadUsage(bot, chatId)
		return
	}
	updates := parseKeyValueArgs(parts[1:])
	product, err := service.UpdateProduct(parts[0], updates)
	if err != nil {
		sendBadUsage(bot, chatId)
		return
	}
	sendText(bot, chatId, fmt.Sprintf(core.STATIC_PRODUCT_UPDATED_TEXT, product.ID, product.Name, product.Price))
	notifyOwner(bot, service, fmt.Sprintf("product updated: %s (%s)", product.Name, product.ID))
}

func HandleDeleteProduct(bot *kappelas.Bot, service *BotService, chatId int64, isAdmin bool, args string) {
	if !isAdmin {
		sendUnauthorized(bot, chatId)
		return
	}
	productID := strings.TrimSpace(args)
	if productID == "" {
		sendBadUsage(bot, chatId)
		return
	}
	if err := service.DeleteProduct(productID); err != nil {
		sendBadUsage(bot, chatId)
		return
	}
	sendText(bot, chatId, fmt.Sprintf(core.STATIC_PRODUCT_DELETED_TEXT, productID))
	notifyOwner(bot, service, fmt.Sprintf("product deleted: %s", productID))
}

func HandleAdminListOrders(bot *kappelas.Bot, service *BotService, chatId int64, isAdmin bool) {
	if !isAdmin {
		sendUnauthorized(bot, chatId)
		return
	}
	orders, err := service.ListOrders()
	if err != nil {
		sendText(bot, chatId, core.STATIC_FALLBACK_TEXT)
		return
	}
	sendOrEdit(bot, service, chatId, "admin-orders", formatOrders(orders), nil)
}

func HandleAdminUpdateOrderStatus(bot *kappelas.Bot, service *BotService, chatId int64, isAdmin bool, args string) {
	if !isAdmin {
		sendUnauthorized(bot, chatId)
		return
	}
	parts := strings.Fields(args)
	if len(parts) != 2 {
		sendBadUsage(bot, chatId)
		return
	}
	order, err := service.UpdateOrderStatus(parts[0], OrderStatus(strings.ToLower(parts[1])))
	if err != nil {
		sendBadUsage(bot, chatId)
		return
	}
	sendText(bot, chatId, fmt.Sprintf(core.STATIC_ORDER_UPDATED_TEXT, order.ID, order.Status))
	notifyOrderChat(bot, order)
	notifyOwner(bot, service, fmt.Sprintf("order %s status updated to %s", order.ID, order.Status))
}

func HandleCommandNotImplemented(bot *kappelas.Bot, chatId int64) {
	ctx := context.Background()
	bot.Messages.Send(ctx, kappelas.SendMessageParams{
		ChatID: chatId,
		Text:   core.STATIC_COMMAND_NOT_IMPLEMENTED_TEXT,
	})
}

func sendText(bot *kappelas.Bot, chatId int64, text string) {
	ctx := context.Background()
	bot.Messages.Send(ctx, kappelas.SendMessageParams{
		ChatID:         chatId,
		Text:           text,
		DeletePrevious: true,
	})
}

func sendOrEdit(bot *kappelas.Bot, service *BotService, chatId int64, key string, text string, markup any) {
	ctx := context.Background()
	if messageID, ok := service.GetMessageRef(chatId, key); ok {
		params := kappelas.EditMessageParams{
			ChatID:    chatId,
			MessageID: messageID,
			NewText:   text,
		}
		if markup != nil {
			if extraData, err := json.Marshal(markup); err == nil {
				params.NewExtraData = extraData
			}
		}
		if _, err := bot.Messages.Edit(ctx, params); err == nil {
			return
		}
	}

	result, err := bot.Messages.Send(ctx, kappelas.SendMessageParams{
		ChatID:         chatId,
		Text:           text,
		ReplyMarkup:    markup,
		DeletePrevious: true,
	})
	if err == nil && result != nil {
		_ = service.SaveMessageRef(chatId, key, result.MessageID)
	}
}

func sendProductImages(bot *kappelas.Bot, chatId int64, product *ShopProduct) {
	if len(product.Images) == 0 {
		return
	}
	cards := make([]kappelas.CarouselCard, 0, len(product.Images))
	for i, imageURL := range product.Images {
		imageURL = strings.TrimSpace(imageURL)
		if imageURL == "" {
			continue
		}
		cards = append(cards, kappelas.CarouselCard{
			ID:       fmt.Sprintf("%s-image-%d", product.ID, i+1),
			Title:    fmt.Sprintf("%s image %d", product.Name, i+1),
			ImageURL: &imageURL,
		})
	}
	if len(cards) == 0 {
		return
	}
	ctx := context.Background()
	_, _ = bot.Messages.SendCarousel(ctx, kappelas.SendCarouselParams{
		ChatID:   chatId,
		Text:     "Product images",
		Carousel: cards,
	})
}

func notifyOrderChat(bot *kappelas.Bot, order *ShopOrder) {
	if order.ChatID == 0 {
		return
	}
	ctx := context.Background()
	bot.Messages.Send(ctx, kappelas.SendMessageParams{
		ChatID: order.ChatID,
		Text:   fmt.Sprintf("Your order %s is now %s.", order.ID, order.Status),
	})
}

func notifyOwner(bot *kappelas.Bot, service *BotService, event string) {
	ownerChatID, ok := service.GetOwnerChat()
	if !ok {
		return
	}
	ctx := context.Background()
	bot.Messages.Send(ctx, kappelas.SendMessageParams{
		ChatID: ownerChatID,
		Text:   fmt.Sprintf(core.STATIC_OWNER_EVENT_TEXT, event),
	})
}

func sendBadUsage(bot *kappelas.Bot, chatId int64) {
	sendText(bot, chatId, core.STATIC_BAD_COMMAND_USAGE_TEXT)
}

func sendUnauthorized(bot *kappelas.Bot, chatId int64) {
	sendText(bot, chatId, core.STATIC_UNAUTHORIZED_TEXT)
}

func splitPipeArgs(args string) []string {
	rawParts := strings.Split(args, "|")
	parts := make([]string, 0, len(rawParts))
	for _, part := range rawParts {
		part = strings.TrimSpace(part)
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

func parseKeyValueArgs(parts []string) map[string]string {
	values := map[string]string{}
	for _, part := range parts {
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		key = strings.ToLower(strings.TrimSpace(key))
		value = strings.TrimSpace(value)
		if key != "" && value != "" {
			values[key] = value
		}
	}
	return values
}

func formatOrders(orders []ShopOrder) string {
	if len(orders) == 0 {
		return core.STATIC_NO_ORDERS_TEXT
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].ID < orders[j].ID
	})
	var b strings.Builder
	b.WriteString("Orders:\n")
	for _, order := range orders {
		b.WriteString(fmt.Sprintf(
			"\nOrder ID: %s\nUser ID: %s\nChat ID: %d\nProduct ID: %s\nQuantity: %d\nTotal: $%.2f\nStatus: %s\n",
			order.ID,
			order.UserID,
			order.ChatID,
			order.ProductID,
			order.Quantity,
			order.TotalPrice,
			order.Status,
		))
	}
	return b.String()
}

func formatCategories(categories []ShopCategory) string {
	if len(categories) == 0 {
		return "No categories found."
	}
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Name < categories[j].Name
	})
	var b strings.Builder
	b.WriteString("Categories:\n")
	for _, category := range categories {
		b.WriteString(fmt.Sprintf("\nCategory ID: %s\nName: %s\n", category.ID, category.Name))
	}
	return b.String()
}
