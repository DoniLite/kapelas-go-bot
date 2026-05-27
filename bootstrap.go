package main

import (
	"log"
	"path"

	"github.com/Arnel7/kappelas-sdk-go"
	"github.com/DoniLite/kapelas-bot/conf"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"os"
)

var (
	bot  *kappelas.Bot
	user *kappelas.User
	router *gin.Engine
)

func Bootstrap() {
	cwd, err := os.Getwd()
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
}

func debugLoader() {

}

func prodLoader() {

}
