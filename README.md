# Kappelas Golang Bot

[Kappelas](https://kappelas.com) is a social media like platform design for creators and web businessman across the world and Africa. The project is focused on the accessibility and full access grant system like telegram...
Some features are furnished for developers like sdk to build automation and bots. This project is designed to help people understand the power of these tools as you can feel free to use it as the base template.

## Features

- Kappelas bot webhook handling with Gin.
- Product categories and paginated product listing.
- Product details with image carousel support.
- Order placement and order status updates.
- Single admin chat activation via `OWNER_ACCESS_TOKEN`.
- Owner/admin event notifications.
- JSON file persistence for products, categories, orders, bot state, and message references.
- Docker multi-stage build with a test target.
- GitHub Actions CI/CD with Docker image publishing only on version tags.

## Requirements

- Go `1.26.3`
- Docker, optional for containerized usage
- A Kappelas bot token
- A Kappelas platform API key, used for user webhook setup

## Environment

Create your local `.env` from the example:

```bash
cp .env.example .env
```

Available variables:

```env
BOT_TOKEN=your_bot_token_here
BOT_DEBUG=true
BOT_VERSION=1.0.0
BOT_IS_DEVELOPMENT=true
BOT_PLATFORM_API_KEY=your_platform_api_key_here
SERVER_PORT=8080
SERVER_HOST=localhost
SERVER_PATH=https://yourdomain.com
OWNER_ACCESS_TOKEN=
```

Notes:

- `BOT_TOKEN` is required for the bot client.
- `BOT_PLATFORM_API_KEY` is used to initialize the Kappelas user client.
- `BOT_IS_DEVELOPMENT=true` starts the bot websocket listener and seeds demo products.
- `BOT_IS_DEVELOPMENT=false` runs the HTTP webhook server.
- `SERVER_PATH` should be your public production base URL, for example `https://bot.example.com`.
- `OWNER_ACCESS_TOKEN` can be left empty. The bot will generate, persist, and log one at startup.

## Admin Chat Setup

The bot does not rely on Kappelas sender IDs for admin detection. Instead, it uses a single admin chat.

1. Start the bot.
2. Copy the token from the startup log:

```txt
Owner access token ready: <token>
```

3. In the chat that should become the admin chat, send:

```txt
/request_owner_access <token>
```

That chat is then stored as the only admin chat. If another chat validates the token later, it replaces the previous admin chat.

Once the admin chat is active, it can retrieve the token with:

```txt
/owner_access_token
```

## Run Locally

Install dependencies and run:

```bash
go mod download
go run .
```

Run tests:

```bash
go test ./...
```

If your Go cache is not writable in your environment:

```bash
GOCACHE=/tmp/go-build go test ./...
```

## Bot Commands

User commands:

```txt
/start
/help
/list_products
/view_product <product_id>
/place_order <product_id> [quantity]
/my_orders
/request_owner_access <token>
```

Admin commands:

```txt
/admin_view_categories
/add_category <name>
/update_category <category_id> | <name>
/delete_category <category_id>
/add_product <category_id> | <name> | <price> | <description> | <image1,image2>
/update_product <product_id> | name=<name> | price=<price> | description=<description> | category=<category_id> | images=<image1,image2>
/delete_product <product_id>
/admin_list_orders
/admin_update_order_status <order_id> <pending|confirmed|shipped|delivered|cancelled>
/owner_access_token
```

Examples:

```txt
/add_category Electronics
/add_product category-id | Smartphone | 699.99 | A high-end smartphone | https://example.com/a.jpg,https://example.com/b.jpg
/update_product product-id | price=649.99 | name=Smartphone Pro
/admin_update_order_status order-id confirmed
```

## Product Images

Seeded demo products include real image URLs. Product details are sent with text and, when images exist, a Kappelas carousel using each product image URL.

The Kappelas SDK `SendPhoto` API expects uploaded file bytes, not remote URLs, so URL-based product images are rendered with `SendCarousel`.

## Persistence

The project uses a JSON store.

- Development mode: `./data_dev`
- Test mode: `/tmp/kappelas-go-bot/data`
- Production mode: the app user's config directory, for Docker:

```txt
/home/app/.config/kappelas-go-bot/data
```

Persisted collections include:

- `products`
- `product_categories`
- `orders`
- `bot_state`

For Docker production deployments, mount a volume to keep data across container restarts:

```bash
docker run --env-file .env \
  -p 8080:8080 \
  -v kappelas-data:/home/app/.config/kappelas-go-bot/data \
  kappelas-go-bot:local
```

## Docker

Build the test target:

```bash
docker build --target test --build-arg GO_VERSION=1.26.3 .
```

Build the runtime image:

```bash
docker build --build-arg GO_VERSION=1.26.3 -t kappelas-go-bot:local .
```

Run the image:

```bash
docker run --env-file .env -p 8080:8080 kappelas-go-bot:local
```

## Webhooks

The bot registers:

```txt
/webhook/bot
/webhook/user
```

In production, `SERVER_PATH` is used as the public base URL. For example:

```env
SERVER_PATH=https://bot.example.com
```

The bot webhook URL becomes:

```txt
https://bot.example.com/webhook/bot
```

## CI/CD

GitHub Actions workflow: `.github/workflows/ci-cd.yml`

On pull requests and pushes to `master`, CI runs:

```bash
go test ./...
docker build --target test --build-arg GO_VERSION=1.26.3 .
```

Docker images are built and published to GHCR only when pushing a Git tag matching `v*`.

Release example:

```bash
git tag v1.0.0
git push origin v1.0.0
```

The published image is tagged from the Git tag and commit SHA.

## Development Notes

- In development mode, `SeedProducts()` seeds demo categories and products with real image URLs.
- Existing `data_dev` records are not automatically overwritten. Clear or update your local JSON data if you want to reseed from scratch.
- Admin notifications are sent to the stored admin chat.
- Customer order status notifications are sent to the chat that created the order.
- Reusable views try to edit existing bot messages first; if editing fails, the bot sends a new message with `DeletePrevious`.
