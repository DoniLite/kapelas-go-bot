package bot

import "strings"

type Command string

func (c Command) Match(input string) bool {
	cmd := ParseCommand(input)
	return cmd == c
}

func (c Command) String() string {
	return string(c)
}

func ParseCommand(input string) Command {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}

	command := strings.Fields(input)[0]
	if idx := strings.Index(command, "?"); idx >= 0 {
		command = command[:idx]
	}
	if !strings.HasPrefix(command, "/") {
		command = "/" + command
	}

	return Command(strings.ToLower(command))
}

func CommandArgs(input string) string {
	parts := strings.Fields(strings.TrimSpace(input))
	if len(parts) < 2 {
		return ""
	}
	return strings.Join(parts[1:], " ")
}

func (c Command) IsKnown() bool {
	for _, command := range AllCommands {
		if c == command {
			return true
		}
	}
	return false
}

const (
	StartCommand                  Command = "/start"
	HelpCommand                   Command = "/help"
	ListProductsCommand           Command = "/list_products"
	ViewProductsCommand           Command = "/view_product"
	PlaceOrderCommand             Command = "/place_order"
	MyOrdersCommand               Command = "/my_orders"
	RequestOwnerAccessCommand     Command = "/request_owner_access"
	AddProductCommand             Command = "/add_product"
	UpdateProductCommand          Command = "/update_product"
	DeleteProductCommand          Command = "/delete_product"
	AdminListOrdersCommand        Command = "/admin_list_orders"
	AdminUpdateOrderStatusCommand Command = "/admin_update_order_status"
	AddCategoryCommand            Command = "/add_category"
	UpdateCategoryCommand         Command = "/update_category"
	DeleteCategoryCommand         Command = "/delete_category"
	AdminViewCategoriesCommand    Command = "/admin_view_categories"
	OwnerAccessTokenCommand       Command = "/owner_access_token"
)

var AllCommands = []Command{
	StartCommand,
	HelpCommand,
	ListProductsCommand,
	ViewProductsCommand,
	PlaceOrderCommand,
	MyOrdersCommand,
	RequestOwnerAccessCommand,
	AddProductCommand,
	UpdateProductCommand,
	DeleteProductCommand,
	AdminListOrdersCommand,
	AdminUpdateOrderStatusCommand,
	AddCategoryCommand,
	UpdateCategoryCommand,
	DeleteCategoryCommand,
	AdminViewCategoriesCommand,
	OwnerAccessTokenCommand,
}
