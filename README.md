# 💱 Deals: Currency Exchange API

A REST API in Go that converts an amount between currencies, validates the request, and saves every exchange to MongoDB. It is built with the [Gin](https://github.com/gin-gonic/gin) web framework and also connects to Redis. MongoDB and Redis run in Docker.

## ✨ Features

- **Currency conversion:** `POST /api/exchange` takes an amount and a currency pair and returns the converted amount.
- **Validation:**
  - the amount must be positive,
  - `created_at` can't be in the future,
  - the currency pair must have an available rate.
- **Saved to MongoDB:** every exchange is stored in the `Currs` collection.
- **Health check:** `GET /api/healthchecker` confirms Redis and the server are up.
- **Config from env:** port and connection strings come from `app.env` (loaded with Viper).

## 🛠️ Built With

| | |
| --- | --- |
| Language | [Go](https://go.dev) |
| Web framework | [Gin](https://github.com/gin-gonic/gin) |
| Database | [MongoDB](https://www.mongodb.com) (official Go driver) |
| Key-value store | [Redis](https://redis.io) (go-redis v8) |
| DevOps | Docker Compose |

## 🗂️ Project Structure

```
task/
├── main.go               # Server setup, MongoDB + Redis connections, routes
├── config/default.go     # Loads app.env
├── models/data.go        # Rate and BasicExchange structs
├── routes/create.go      # Saves an exchange to MongoDB
├── databases/database.go # MongoDB client
├── Collection/           # Collection helper
├── docker-compose.yml    # MongoDB + Redis
└── app.env               # Port and connection settings
```

## 🔌 API

### `POST /api/exchange`

**Request**

```json
{
  "rates": [{ "cfrom": "USD", "cto": "GBP", "conv": 1.5 }],
  "amount": 5000,
  "codeId": 1,
  "created_at": "2022-10-14T17:31:03.735Z"
}
```

**Response `201 Created`**

```json
{
  "success": true,
  "message": "data created successfully",
  "data": {
    "rates": [{ "cfrom": "USD", "cto": "GBP", "conv": 1.5 }],
    "codeId": 1,
    "_id": "63499d39985b28b42b43c263",
    "created_at": "2022-10-14T17:31:03.735Z",
    "amount": 5000,
    "NewAmount": 7500
  }
}
```

**Errors `400 Bad Request`**

| Case | Message |
| --- | --- |
| `created_at` is in the future | `Sorry but you cannot do the exchange for the currencies the time of createdAt is in future` |
| No rate for the currency pair, or amount is 0 or less | `No Rate Available` |

### `GET /api/healthchecker`

```json
{ "status": "success", "message": "Welcome to Golang with Redis and MongoDB" }
```

## 🚀 Getting Started

**Requirements:** Go and Docker.

```bash
git clone https://github.com/YoumnaSalloum/Deals_project.git
cd Deals_project/task
docker-compose up -d    # starts MongoDB (host port 6000) and Redis (6379)
go run main.go          # API on http://localhost:8000
```

Postman documentation: [gin-server-doc](https://documenter.getpostman.com/view/8916756/SztG3mKJ?version=latest)

## 👩‍💻 Author

**Youmna Salloum** · [Portfolio](https://youmnasalloum.github.io/) · [LinkedIn](https://www.linkedin.com/in/youmna-salloum) · [GitHub](https://github.com/YoumnaSalloum)
