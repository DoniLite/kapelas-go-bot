package bot

import (
	"context"
	"fmt"
	"log"

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
	bot.Messages.Send(ctx, kappelas.SendMessageParams{
		ChatID:      chatId,
		Text:        fmt.Sprintf(core.STATIC_WELCOME_TEXT, senderName, "doni"),
		ReplyMarkup: markup,
	})
}

func HandleHelpCommand(bot *kappelas.Bot, chatId int64) {
	ctx := context.Background()
	bot.Messages.Send(ctx, kappelas.SendMessageParams{
		ChatID: chatId,
		Text:   core.STATIC_HELP_TEXT,
	})
}

func HandleListProducts(bot *kappelas.Bot, service *BotService, chatId int64, query map[string]string) {
	ctx := context.Background()
	categoryId := query["category"]
	log.Printf("Listing products for category id: %s", categoryId)
	category, err := service.GetCategoryById(categoryId)
	if err != nil {
		bot.Messages.Send(ctx, kappelas.SendMessageParams{
			ChatID: chatId,
			Text:   core.STATIC_FALLBACK_TEXT,
		})
		return
	}
	markup, err := service.ListAvailableProducts(categoryId)
	if err != nil {
		bot.Messages.Send(ctx, kappelas.SendMessageParams{
			ChatID: chatId,
			Text:   core.STATIC_FALLBACK_TEXT,
		})
		return
	}
	bot.Messages.Send(ctx, kappelas.SendMessageParams{
		ChatID:      chatId,
		Text:        fmt.Sprintf(core.STATIC_LIST_PRODUCTS_BY_CATEGORY_TEXT, category.Name),
		ReplyMarkup: markup,
	})
}
