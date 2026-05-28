package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
	deps.Bot.OnMessage(func(m *kappelas.Message) {
		if bot.StartCommand.Match(*m.Text) {
			bot.HandleStartCommand(deps.Bot, botService, m.ChatID, *m.SenderName)
			return
		}
		if bot.HelpCommand.Match(*m.Text) {
			bot.HandleHelpCommand(deps.Bot, m.ChatID)
			return
		}
		if bot.ListProductsCommand.Match(*m.Text) {
			bot.HandleListProducts(deps.Bot, botService, m.ChatID, map[string]string{})
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
		log.Printf("Received callback query: command=%s, userId=%s, productId=%s, categoryId=%s, orderId=%s",
			query.Command, query.UserId, query.ProductId, query.CategoryId, query.OrderId)
		switch query.Command {
		case bot.StartCommand:
			{
				bot.HandleStartCommand(deps.Bot, botService, cq.ChatID, *cq.SenderUsername)
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
				})
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

	b := &BotCallbackQuery{
		Command:    bot.Command(cmdPart),
		UserId:     getFirst("user", "user_id", "userid", "uid"),
		ProductId:  getFirst("product", "product_id", "productid"),
		CategoryId: getFirst("category", "category_id", "categoryid"),
		OrderId:    getFirst("order", "order_id", "orderid"),
	}

	return b, nil
}
