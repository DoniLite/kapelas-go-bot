package main

import (
	"context"
	"log"

	"github.com/Arnel7/kappelas-sdk-go"
	"github.com/DoniLite/kapelas-bot/conf"
	"github.com/DoniLite/kapelas-bot/handlers"
	botLib "github.com/DoniLite/kapelas-bot/handlers/bot"
	"github.com/gin-gonic/gin"
)

var (
	bot    *kappelas.Bot
	user   *kappelas.User
	router *gin.Engine
)

func Bootstrap() {
	ctx := context.Background()

	env := conf.GetEnv()

	bot = kappelas.NewBot(env.GetString(conf.BOT_TOKEN))
	user = kappelas.NewUser(env.GetString(conf.BOT_PLATFORM_API_KEY))
	router = gin.Default()

	bot.Webhooks.Set(ctx, kappelas.SetWebhookParams{
		URL: env.GetWebHookBotURL(),
	})
	user.Webhooks.Set(ctx, kappelas.SetWebhookParams{
		URL: env.GetWebHookUserURL(),
	})
	log.Printf("Bot webhook set to: %s", env.GetWebHookBotURL())
	log.Printf("User webhook set to: %s", env.GetWebHookUserURL())

	router = handlers.BuildRouter(&handlers.RouterDeps{
		Router: router,
		Bot:    bot,
		User:   user,
	})

	isDev := env.GetBool(conf.BOT_IS_DEVELOPMENT)
	if isDev {
		log.Printf("Running in dev mode, starting the channel to listen for bot updates...")
		botLib.SeedProducts()
		bot.Start()
		select {}
	}
	router.Run(":" + env.GetString(conf.SERVER_PORT))
}
