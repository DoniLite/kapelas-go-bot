package handlers

import (
	"github.com/Arnel7/kappelas-sdk-go"
	"github.com/DoniLite/kapelas-bot/handlers/bot"
	"github.com/gin-gonic/gin"
)

type RouterDeps struct {
	router *gin.Engine
	bot *kappelas.Bot
	user *kappelas.User
}

func BuildRouter(deps *RouterDeps) *gin.Engine {
	botService := bot.NewBotService(deps.bot)
	return deps.router
}