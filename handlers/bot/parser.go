package bot

import (
	"github.com/Arnel7/kappelas-sdk-go"
)

func TotalizePrice(price float64, quantity int) float64 {
	return price * float64(quantity)
}

func GenerateProductListMarker(products ...ShopProduct) *kappelas.InlineKeyboard {
	var buttons [][]kappelas.InlineKeyboardButton
	for _, product := range products {
		buttons = append(buttons, []kappelas.InlineKeyboardButton{
			{
				Text:         product.Name,
				CallbackData: new(product.ID),
			},
		})
	}
	return &kappelas.InlineKeyboard{
		InlineKeyboard: buttons,
	}
}
