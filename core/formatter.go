package core

var (
	STATIC_WELCOME_TEXT = `
👋🏽 Hello %s! Welcome to *My Kappelas Shop*.

Bellow you have all my products listed by categories.
To get started, just click on a category and I'll show you the products I have in that category.

If you have any questions, just ask me here. I'm always happy to help!
You can contact me on this username: @%s

Happy shopping! 🛍️
`

	STATIC_FALLBACK_TEXT = `
Sorry, I didn't understand that. Please try again.
Or type /help to see the available commands.
`

	STATIC_HELP_TEXT = `Here are the available commands:

/start - Start the bot and see the welcome message
/help - Show this help message
/list_products - Choose a category and browse products
/view_product <product_id> - View product details
/place_order <product_id> [quantity] - Place an order
/my_orders - View your orders
/request_owner_access <token> - Request owner access

Admin commands:
/admin_view_categories - List categories
/add_category <name> - Add a category
/update_category <category_id> | <name> - Update a category
/delete_category <category_id> - Delete a category
/add_product <category_id> | <name> | <price> | <description> | <image1,image2>
/update_product <product_id> | name=<name> | price=<price> | description=<description> | category=<category_id> | images=<image1,image2>
/delete_product <product_id> - Delete a product
/admin_list_orders - List all orders
/admin_update_order_status <order_id> <pending|confirmed|shipped|delivered|cancelled>
/owner_access_token - Show the seeded owner access token

You can also click on the buttons to navigate through the categories and products.
`

	STATIC_LIST_PRODUCTS_BY_CATEGORY_TEXT = `Here are the products in the %s category:`

	STATIC_LIST_PRODUCTS_BY_CATEGORY_PAGE_TEXT = `Here are the products in the %s category:

Page %d/%d`

	STATIC_EMPTY_PRODUCT_CATEGORY_TEXT = `No products are available in the %s category yet.`

	STATIC_CHOOSE_CATEGORY_TEXT = `Choose a category to browse products:`

	STATIC_VIEW_PRODUCT_TEXT = `%s

%s

Price: $%.2f

ID: %s
Category: %s`

	STATIC_COMMAND_NOT_IMPLEMENTED_TEXT = `This command is not available yet.`

	STATIC_BAD_COMMAND_USAGE_TEXT = `Invalid command usage.

Type /help to see the expected format.`

	STATIC_UNAUTHORIZED_TEXT = `You are not allowed to use this command.`

	STATIC_ORDER_CREATED_TEXT = `Order created successfully.

Order: %s
Product: %s
Quantity: %d
Total: $%.2f
Status: %s`

	STATIC_ORDER_UPDATED_TEXT = `Order %s status updated to %s.`

	STATIC_NO_ORDERS_TEXT = `No orders found.`

	STATIC_PRODUCT_CREATED_TEXT = `Product created successfully.

ID: %s
Name: %s
Price: $%.2f`

	STATIC_PRODUCT_UPDATED_TEXT = `Product updated successfully.

ID: %s
Name: %s
Price: $%.2f`

	STATIC_PRODUCT_DELETED_TEXT = `Product %s deleted successfully.`

	STATIC_CATEGORY_CREATED_TEXT = `Category created successfully.

ID: %s
Name: %s`

	STATIC_CATEGORY_UPDATED_TEXT = `Category updated successfully.

ID: %s
Name: %s`

	STATIC_CATEGORY_DELETED_TEXT = `Category %s deleted successfully.`

	STATIC_OWNER_ACCESS_GRANTED_TEXT = `Owner access token accepted.

This chat is now the admin chat.

Chat ID: %d`

	STATIC_OWNER_ACCESS_TOKEN_TEXT = `Owner access token:

%s`

	STATIC_OWNER_EVENT_TEXT = `Owner event: %s`
)
