package core

import (
	"strconv"

	"github.com/Arnel7/kappelas-sdk-go"
)

func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

type KappelasMessage = kappelas.Message
type KappelasCallbackQuery = kappelas.CallbackQuery