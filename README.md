# Go Contributors API

A GraphQL API built with Go for managing contributors, featuring MongoDB for data storage and Redis for caching and rate limiting.

## Features

- GraphQL API using gqlgen
- MongoDB for persistent storage
- Redis for caching and rate limiting
- Clean architecture with separation of concerns
- Comprehensive test coverage

## Prerequisites

- Go 1.21 or higher
- MongoDB
- Redis
- Make (optional, for using Makefile)

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/aakashdeepsil/go-contributors-api.git
2. Install dependencies
   ```bash
   go mod download
3. Copy the example environment file:
   ```bash
   cp .env.example .env
4. Update the environment variables in .env as needed
5. Run the server:
   ```bash
   go run cmd/server/main.go

## Project Structure

- cmd/: Application entry points
- internal/: Private application code
   - config/: Configuration management
   - database/: Database connections and repositories
   - graph/: GraphQL schema and resolvers
   - middleware/: HTTP middleware
   - service/: Business logic
- pkg/: Public libraries that can be used by external applications
