# Milkly Member Service

A Go-based REST API service for member management using Fiber framework and MongoDB.

## Features

- RESTful API for member management
- MongoDB integration
- Swagger documentation
- Dependency injection with Uber Dig
- Input validation
- Health check endpoint

## Prerequisites

- Go 1.25 or later
- MongoDB (local or remote)
- Make (optional)

## Setup

1. Clone the repository
2. Copy environment file:
   ```bash
   cp .env.example .env
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Start MongoDB (if running locally)
5. Run the application:
   ```bash
   make run
   # or
   go run cmd/milkly-member/main.go
   ```

## API Endpoints

### Member Management
- `POST /api/v1/members` - Create a new member
- `GET /api/v1/members` - List all members (with pagination)
- `GET /api/v1/members/{id}` - Get a specific member
- `PUT /api/v1/members/{id}` - Update a member
- `DELETE /api/v1/members/{id}` - Delete a member

### Health Check
- `GET /health` - Service health check

### Documentation
- `GET /swagger/index.html` - Swagger UI

## Configuration

Environment variables:
- `PORT` - Server port (default: 8080)
- `MONGO_URL` - MongoDB connection URL (default: mongodb://localhost:27017)
- `DB_NAME` - Database name (default: milkly-member)

## Project Structure

```
.
├── cmd/
│   └── milkly-member/
│       └── main.go          # Application entry point
├── config/
│   └── config.go           # Configuration management
├── controller/
│   ├── controllers.go      # Router and middleware setup
│   ├── member_controller.go # Member API handlers
│   └── health_controller.go # Health check handler
├── di/
│   └── container.go        # Dependency injection setup
├── docs/
│   └── docs.go            # Swagger documentation
├── entity/
│   └── member.go          # Data models and DTOs
├── repository/
│   └── member_repository.go # Database layer
├── service/
│   └── member_service.go   # Business logic layer
├── .env.example           # Environment variables template
├── go.mod                 # Go module file
└── README.md             # This file
```

## Development

### Generate Swagger Documentation
```bash
swag init -g cmd/milkly-member/main.go -o docs/
```

### Run Tests
```bash
go test ./...
```

### Build
```bash
go build -o bin/milkly-member cmd/milkly-member/main.go
```

## API Usage Examples

### Create a Member
```bash
curl -X POST http://localhost:8080/api/v1/members \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890"
  }'
```

### Get All Members
```bash
curl http://localhost:8080/api/v1/members?limit=10&offset=0
```

### Get a Member
```bash
curl http://localhost:8080/api/v1/members/{member_id}
```

### Update a Member
```bash
curl -X PUT http://localhost:8080/api/v1/members/{member_id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "status": "active"
  }'
```

### Delete a Member
```bash
curl -X DELETE http://localhost:8080/api/v1/members/{member_id}
```