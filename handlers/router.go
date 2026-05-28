package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Arnel7/kappelas-sdk-go"
	"github.com/DoniLite/kapelas-bot/conf"
	"github.com/DoniLite/kapelas-bot/core"
	"github.com/DoniLite/kapelas-bot/handlers/bot"
	"github.com/gin-gonic/gin"
)

type RouterDeps struct {
	Router *gin.Engine
	Bot    *kappelas.Bot
	User   *kappelas.User
}

type BotCallbackQuery struct {
	Command    bot.Command
	UserId     string
	ProductId  string
	CategoryId string
	OrderId    string
	Page       int
	Quantity   int
}

func HandleBotWebHook(deps *RouterDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		deps.Bot.HandleWebhook(body)
		c.Writer.WriteHeader(http.StatusOK)
	}
}

func BuildRouter(deps *RouterDeps) *gin.Engine {
	ctx := context.Background()
	botService := bot.NewBotService(deps.Bot, bot.NewBotRepository())
	if token, err := botService.EnsureOwnerAccessToken(); err == nil {
		log.Printf("Owner access token ready: %s", token)
	}
	deps.Bot.OnMessage(func(m *kappelas.Message) {
		if m.Text == nil {
			deps.Bot.Messages.Send(ctx, kappelas.SendMessageParams{
				ChatID: m.ChatID,
				Text:   core.STATIC_FALLBACK_TEXT,
			})
			return
		}

		command := bot.ParseCommand(*m.Text)
		args := bot.CommandArgs(*m.Text)
		userID := userIDFromMessage(m)
		isAdmin := botService.IsOwnerChat(m.ChatID)
		switch command {
		case bot.StartCommand:
			bot.HandleStartCommand(deps.Bot, botService, m.ChatID, safeString(m.SenderName))
			return
		case bot.HelpCommand:
			bot.HandleHelpCommand(deps.Bot, m.ChatID)
			return
		case bot.ListProductsCommand:
			bot.HandleListProducts(deps.Bot, botService, m.ChatID, map[string]string{})
			return
		case bot.ViewProductsCommand:
			bot.HandleViewProduct(deps.Bot, botService, m.ChatID, map[string]string{
				"product": args,
			})
			return
		case bot.PlaceOrderCommand:
			productID, quantity := parsePlaceOrderArgs(args)
			bot.HandlePlaceOrder(deps.Bot, botService, m.ChatID, userID, map[string]string{
				"product":  productID,
				"quantity": strconv.Itoa(quantity),
			})
			return
		case bot.MyOrdersCommand:
			bot.HandleMyOrders(deps.Bot, botService, m.ChatID, userID)
			return
		case bot.RequestOwnerAccessCommand:
			bot.HandleRequestOwnerAccess(deps.Bot, botService, m.ChatID, args)
			return
		case bot.OwnerAccessTokenCommand:
			bot.HandleOwnerAccessToken(deps.Bot, botService, m.ChatID, isAdmin)
			return
		case bot.AdminViewCategoriesCommand:
			bot.HandleAdminViewCategories(deps.Bot, botService, m.ChatID, isAdmin)
			return
		case bot.AddCategoryCommand:
			bot.HandleAddCategory(deps.Bot, botService, m.ChatID, isAdmin, args)
			return
		case bot.UpdateCategoryCommand:
			bot.HandleUpdateCategory(deps.Bot, botService, m.ChatID, isAdmin, args)
			return
		case bot.DeleteCategoryCommand:
			bot.HandleDeleteCategory(deps.Bot, botService, m.ChatID, isAdmin, args)
			return
		case bot.AddProductCommand:
			bot.HandleAddProduct(deps.Bot, botService, m.ChatID, isAdmin, args)
			return
		case bot.UpdateProductCommand:
			bot.HandleUpdateProduct(deps.Bot, botService, m.ChatID, isAdmin, args)
			return
		case bot.DeleteProductCommand:
			bot.HandleDeleteProduct(deps.Bot, botService, m.ChatID, isAdmin, args)
			return
		case bot.AdminListOrdersCommand:
			bot.HandleAdminListOrders(deps.Bot, botService, m.ChatID, isAdmin)
			return
		case bot.AdminUpdateOrderStatusCommand:
			bot.HandleAdminUpdateOrderStatus(deps.Bot, botService, m.ChatID, isAdmin, args)
			return
		}

		deps.Bot.Messages.Send(ctx, kappelas.SendMessageParams{
			ChatID: m.ChatID,
			Text:   core.STATIC_FALLBACK_TEXT,
		})
	})

	deps.Bot.OnCallbackQuery(func(cq *kappelas.CallbackQuery) {
		query, err := parseCallbackQueryData(cq.CallbackData)
		if err != nil {
			deps.Bot.Messages.Send(ctx, kappelas.SendMessageParams{
				ChatID: cq.ChatID,
				Text:   core.STATIC_FALLBACK_TEXT,
			})
			return
		}
		log.Printf("Received callback query: command=%s, userId=%s, productId=%s, categoryId=%s, orderId=%s, page=%d",
			query.Command, query.UserId, query.ProductId, query.CategoryId, query.OrderId, query.Page)
		isAdmin := botService.IsOwnerChat(cq.ChatID)
		userID := userIDFromCallback(cq)
		switch query.Command {
		case bot.StartCommand:
			{
				bot.HandleStartCommand(deps.Bot, botService, cq.ChatID, safeString(cq.SenderUsername))
				break
			}
		case bot.HelpCommand:
			{
				bot.HandleHelpCommand(deps.Bot, cq.ChatID)
				break
			}
		case bot.ListProductsCommand:
			{
				bot.HandleListProducts(deps.Bot, botService, cq.ChatID, map[string]string{
					"category": query.CategoryId,
					"page":     strconv.Itoa(query.Page),
				})
				break
			}
		case bot.ViewProductsCommand:
			{
				bot.HandleViewProduct(deps.Bot, botService, cq.ChatID, map[string]string{
					"product": query.ProductId,
				})
				break
			}
		case bot.PlaceOrderCommand:
			{
				bot.HandlePlaceOrder(deps.Bot, botService, cq.ChatID, userID, map[string]string{
					"product":  query.ProductId,
					"quantity": strconv.Itoa(query.Quantity),
				})
				break
			}
		case bot.MyOrdersCommand:
			{
				bot.HandleMyOrders(deps.Bot, botService, cq.ChatID, userID)
				break
			}
		case bot.AdminViewCategoriesCommand:
			{
				bot.HandleAdminViewCategories(deps.Bot, botService, cq.ChatID, isAdmin)
				break
			}
		case bot.AdminListOrdersCommand:
			{
				bot.HandleAdminListOrders(deps.Bot, botService, cq.ChatID, isAdmin)
				break
			}
		case bot.OwnerAccessTokenCommand:
			{
				bot.HandleOwnerAccessToken(deps.Bot, botService, cq.ChatID, isAdmin)
				break
			}
		default:
			deps.Bot.Messages.Send(ctx, kappelas.SendMessageParams{
				ChatID: cq.ChatID,
				Text:   core.STATIC_FALLBACK_TEXT,
			})
		}
	})

	deps.Router.POST(conf.WEBHOOK_BOT_PATH, HandleBotWebHook(deps))
	return deps.Router
}

func parseCallbackQueryData(dataAsURLParts string) (*BotCallbackQuery, error) {
	if strings.TrimSpace(dataAsURLParts) == "" {
		return nil, fmt.Errorf("empty callback data")
	}

	raw := dataAsURLParts
	// split command and query
	var cmdPart string
	var queryPart string
	if idx := strings.Index(raw, "?"); idx >= 0 {
		cmdPart = raw[:idx]
		queryPart = raw[idx+1:]
	} else {
		cmdPart = raw
		queryPart = ""
	}

	if cmdPart == "" {
		return nil, fmt.Errorf("empty command in callback data")
	}
	// ensure leading slash for command consistency
	if !strings.HasPrefix(cmdPart, "/") {
		cmdPart = "/" + cmdPart
	}
	// unescape and normalize command
	if up, err := url.PathUnescape(cmdPart); err == nil {
		cmdPart = up
	}
	cmdPart = strings.ToLower(strings.TrimSpace(cmdPart))

	qvals, _ := url.ParseQuery(queryPart)

	getFirst := func(keys ...string) string {
		for _, k := range keys {
			if v := qvals.Get(k); v != "" {
				return v
			}
		}
		return ""
	}

	page, err := strconv.Atoi(getFirst("page", "p"))
	if err != nil || page < 1 {
		page = 1
	}
	quantity, err := strconv.Atoi(getFirst("quantity", "qty", "q"))
	if err != nil || quantity < 1 {
		quantity = 1
	}

	b := &BotCallbackQuery{
		Command:    bot.Command(cmdPart),
		UserId:     getFirst("user", "user_id", "userid", "uid"),
		ProductId:  getFirst("product", "product_id", "productid"),
		CategoryId: getFirst("category", "category_id", "categoryid"),
		OrderId:    getFirst("order", "order_id", "orderid"),
		Page:       page,
		Quantity:   quantity,
	}

	return b, nil
}

func safeString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func parsePlaceOrderArgs(args string) (string, int) {
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return "", 1
	}
	quantity := 1
	if len(parts) > 1 {
		if value, err := strconv.Atoi(parts[1]); err == nil && value > 0 {
			quantity = value
		}
	}
	return parts[0], quantity
}

func userIDFromMessage(m *kappelas.Message) string {
	if m.SenderID != nil && strings.TrimSpace(*m.SenderID) != "" {
		return *m.SenderID
	}
	return chatUserID(m.ChatID)
}

func userIDFromCallback(cq *kappelas.CallbackQuery) string {
	if strings.TrimSpace(cq.SenderID) != "" {
		return cq.SenderID
	}
	return chatUserID(cq.ChatID)
}

func chatUserID(chatID int64) string {
	return fmt.Sprintf("chat:%d", chatID)
}
