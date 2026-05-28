package bot

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Arnel7/kappelas-sdk-go"
)

func TotalizePrice(price float64, quantity int) float64 {
	return price * float64(quantity)
}

func GenerateProductListMarker(products ...ShopProduct) *kappelas.ReplyKeyboard {
	var buttons [][]string
	for _, product := range products {
		buttons = append(buttons, []string{
			fmt.Sprintf("%s - $%.2f", product.Name, product.Price),
		})
	}
	return &kappelas.ReplyKeyboard{
		Keyboard: buttons,
	}
}

func GenerateProductCategoryListMarker(categories ...ShopCategory) *kappelas.InlineKeyboard {
	var buttons [][]kappelas.InlineKeyboardButton
	for _, category := range categories {
		cb := BuildCallbackData(ListProductsCommand, map[string]string{"category": category.ID})
		buttons = append(buttons, []kappelas.InlineKeyboardButton{
			{
				Text:         category.Name,
				CallbackData: &cb,
			},
		})
	}
	return &kappelas.InlineKeyboard{
		InlineKeyboard: buttons,
	}
}

// BuildCallbackData builds a callback data string from a command and key/value params.
// Example: BuildCallbackData(bot.StartCommand, map[string]string{"category":"c1"}) -> "/start?category=c1"
func BuildCallbackData(cmd Command, params map[string]string) string {
	base := string(cmd)
	if !strings.HasPrefix(base, "/") {
		base = "/" + base
	}
	if len(params) == 0 {
		return base
	}
	vals := url.Values{}
	for k, v := range params {
		vals.Set(k, v)
	}
	return base + "?" + vals.Encode()
}
