package bot

import "github.com/Arnel7/kappelas-sdk-go"

type BotService struct {
	bot        *kappelas.Bot
	repository *BotRepository
}

func NewBotService(bot *kappelas.Bot, repository *BotRepository) *BotService {
	return &BotService{
		bot:        bot,
		repository: repository,
	}
}

func (s *BotService) ListAllAvailableCategories() (*kappelas.InlineKeyboard, error) {
	categories, err := s.repository.GetCategories()
	if err != nil {
		return nil, err
	}
	return GenerateProductCategoryListMarker(categories...), nil
}

func (s *BotService) ListAvailableProducts(categoryIds ...string) (*kappelas.ReplyKeyboard, error) {
	products, err := s.repository.GetProductList(categoryIds...)
	if err != nil {
		return nil, err
	}
	return GenerateProductListMarker(products...), nil
}

func (s *BotService) GetCategoryById(categoryId string) (*ShopCategory, error) {
	return s.repository.GetProductCategory(categoryId)
}
