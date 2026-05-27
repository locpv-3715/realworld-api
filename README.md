# RealWorld API

A fully-featured **Conduit** RealWorld backend API built with **Go**, **Gin**, **GORM**, and **MySQL**.

> Spec: [RealWorld](https://github.com/gothinkster/realworld) — _"The mother of all demo apps"_

---

## Tech Stack

| Layer       | Technology                        |
|-------------|-----------------------------------|
| Language    | Go 1.26+                          |
| Framework   | [Gin](https://github.com/gin-gonic/gin) |
| ORM         | [GORM](https://gorm.io) + MySQL   |
| Auth        | JWT (`golang-jwt/jwt v5`)         |
| Docs        | Swagger (swaggo)                  |
| Config      | `godotenv`                        |

---

## Project Structure

```
.
├── cmd/api/            # Application entry point
├── config/             # Database connection
├── controllers/        # HTTP handlers
├── dto/                # Request / Response data transfer objects
├── middlewares/        # JWT auth middleware
├── models/             # GORM models (User, Article, Comment, etc.)
├── repositories/       # Database access layer
├── routes/             # Route registration
├── services/           # Business logic
├── testhelpers/        # Shared test utilities
├── tests/              # Integration tests
└── utils/              # Utility helpers (slug, etc.)
```

---

## Prerequisites

- Go 1.21+
- MySQL 8+

---

## Setup

### 1. Clone the repository

```bash
git clone https://github.com/your-username/realworld-api.git
cd realworld-api
```

### 2. Configure environment variables

Create a `.env` file in the project root:

```env
DB_DSN=user:password@tcp(127.0.0.1:3306)/realworld?charset=utf8mb4&parseTime=True&loc=Local
JWT_SECRET=your_jwt_secret
PORT=8080
```

### 3. Install dependencies

```bash
go mod download
```

---

## Running the Application

```bash
go run ./cmd/api
```

The server will start at `http://localhost:8080`.

---

## API Documentation (Swagger)

Once the server is running, open:

```
http://localhost:8080/swagger/index.html
```

### Authentication

Pass the JWT token in the `Authorization` header:

```
Authorization: Token <your_jwt_token>
```

---

## API Endpoints

### Users & Authentication

| Method | Endpoint           | Auth     | Description          |
|--------|--------------------|----------|----------------------|
| POST   | `/api/users`       | No       | Register             |
| POST   | `/api/users/login` | No       | Login                |
| GET    | `/api/user`        | Required | Get current user     |
| PUT    | `/api/user`        | Required | Update current user  |

### Profiles

| Method | Endpoint                       | Auth     | Description    |
|--------|--------------------------------|----------|----------------|
| GET    | `/api/profiles/:username`      | Optional | Get profile    |
| POST   | `/api/profiles/:username/follow` | Required | Follow user  |
| DELETE | `/api/profiles/:username/follow` | Required | Unfollow user|

### Articles

| Method | Endpoint                   | Auth     | Description         |
|--------|----------------------------|----------|---------------------|
| GET    | `/api/articles`            | Optional | List articles       |
| GET    | `/api/articles/feed`       | Required | Get feed            |
| GET    | `/api/articles/:slug`      | Optional | Get article         |
| POST   | `/api/articles`            | Required | Create article      |
| PUT    | `/api/articles/:slug`      | Required | Update article      |
| DELETE | `/api/articles/:slug`      | Required | Delete article      |

### Comments

| Method | Endpoint                              | Auth     | Description     |
|--------|---------------------------------------|----------|-----------------|
| GET    | `/api/articles/:slug/comments`        | Optional | Get comments    |
| POST   | `/api/articles/:slug/comments`        | Required | Add comment     |
| DELETE | `/api/articles/:slug/comments/:id`    | Required | Delete comment  |

### Favorites

| Method | Endpoint                           | Auth     | Description         |
|--------|------------------------------------|----------|---------------------|
| POST   | `/api/articles/:slug/favorite`     | Required | Favorite article    |
| DELETE | `/api/articles/:slug/favorite`     | Required | Unfavorite article  |

### Tags

| Method | Endpoint    | Auth | Description   |
|--------|-------------|------|---------------|
| GET    | `/api/tags` | No   | Get all tags  |

---

## Running Tests

### Unit tests

```bash
go test ./...
```

### Unit tests with verbose output

```bash
go test -v ./...
```

### Integration tests

```bash
go test -v ./tests/...
```

---

## Generate Swagger Docs

Install `swag` CLI if not already installed:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Regenerate docs:

```bash
swag init -g cmd/api/main.go
```

---

## License

MIT
