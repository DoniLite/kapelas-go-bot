package main

import (
	"context"
	"log"
	"path"

	"github.com/Arnel7/kappelas-sdk-go"
	"github.com/DoniLite/kapelas-bot/conf"
	"github.com/DoniLite/kapelas-bot/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"os"
)

var (
	bot    *kappelas.Bot
	user   *kappelas.User
	router *gin.Engine
)

func Bootstrap() {
	cwd, err := os.Getwd()
	ctx := context.Background()
	if err != nil {
		log.Fatal("Error getting current working directory, make sure to run the bot from the project root directory")
	}
	DEFAULT_ENV_PATH := path.Join(cwd, "./.env")
	// add your custom env file here to override default env variables
	err = godotenv.Load(DEFAULT_ENV_PATH)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
		go func() {
			bot.Start()
			select {}
		}()
	}
	router.Run(":" + env.GetString(conf.SERVER_PORT))
}
