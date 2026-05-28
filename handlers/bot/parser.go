package bot

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Arnel7/kappelas-sdk-go"
)

const ProductListPageSize = 5

type ProductListPage struct {
	Products   []ShopProduct
	CategoryID string
	Page       int
	PageSize   int
	Total      int
	TotalPages int
}

func TotalizePrice(price float64, quantity int) float64 {
	return price * float64(quantity)
}

func GenerateProductListMarker(page ProductListPage) *kappelas.InlineKeyboard {
	var buttons [][]kappelas.InlineKeyboardButton
	for _, product := range page.Products {
		cb := BuildCallbackData(ViewProductsCommand, map[string]string{
			"product": product.ID,
		})
		buttons = append(buttons, []kappelas.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%s - $%.2f", product.Name, product.Price),
				CallbackData: &cb,
			},
		})
	}

	var paginationButtons []kappelas.InlineKeyboardButton
	if page.Page > 1 {
		cb := BuildCallbackData(ListProductsCommand, map[string]string{
			"category": page.CategoryID,
			"page":     strconv.Itoa(page.Page - 1),
		})
		paginationButtons = append(paginationButtons, kappelas.InlineKeyboardButton{
			Text:         "Previous",
			CallbackData: &cb,
		})
	}
	if page.TotalPages > 1 {
		cb := BuildCallbackData(ListProductsCommand, map[string]string{
			"category": page.CategoryID,
			"page":     strconv.Itoa(page.Page),
		})
		paginationButtons = append(paginationButtons, kappelas.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d/%d", page.Page, page.TotalPages),
			CallbackData: &cb,
		})
	}
	if page.Page < page.TotalPages {
		cb := BuildCallbackData(ListProductsCommand, map[string]string{
			"category": page.CategoryID,
			"page":     strconv.Itoa(page.Page + 1),
		})
		paginationButtons = append(paginationButtons, kappelas.InlineKeyboardButton{
			Text:         "Next",
			CallbackData: &cb,
		})
	}
	if len(paginationButtons) > 0 {
		buttons = append(buttons, paginationButtons)
	}

	return &kappelas.InlineKeyboard{
		InlineKeyboard: buttons,
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
