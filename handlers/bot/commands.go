package bot

import "strings"

type Command string

func (c Command) Match(input string) bool {
	return strings.HasPrefix(input, string(c))
}

const (
	StartCommand       Command = "/start"
	HelpCommand        Command = "/help"
	ListProductsCommand Command = "/list_products"
)
