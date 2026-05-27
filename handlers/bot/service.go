package bot

import "github.com/Arnel7/kappelas-sdk-go"

type BotService struct {
	bot *kappelas.Bot
}

func NewBotService(bot *kappelas.Bot) *BotService {
	return &BotService{
		bot: bot,
	}
}

