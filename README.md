# Alchemist CLI

[![Go Version](https://img.shields.io/badge/go-1.24.4-blue.svg)](https://golang.org/dl/)
[![Release](https://img.shields.io/github/v/release/go-alchemist/alchemist)](https://github.com/go-alchemist/alchemist/releases)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

**Alchemist** is a powerful CLI tool designed to streamline Go project development by providing scaffolding capabilities and database migration management. It helps developers quickly generate boilerplate code and manage database schema changes efficiently.

## ✨ Features

- 🏗️ **Code Scaffolding**: Generate models, handlers, and other Go components
- 🗃️ **Database Migrations**: Create and manage SQL database migrations
- ⚙️ **Flexible Configuration**: Support for environment variables and configuration files
- 🎨 **Colored Output**: Beautiful colored terminal output for better readability
- 🚀 **Cross-platform**: Available for Linux, macOS, and Windows
- 📦 **Zero Dependencies**: Single binary with no external dependencies

## 📥 Installation

### Download Pre-built Binaries

Download the latest release from the [releases page](https://github.com/go-alchemist/alchemist/releases).

### Install from Source

```bash
# Clone the repository
git clone https://github.com/go-alchemist/alchemist.git
cd alchemist

# Build the binary
go build -o alchemist cmd/main.go

# Move to your PATH (optional)
sudo mv alchemist /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/go-alchemist/alchemist/cmd@latest
```

## 🚀 Quick Start

```bash
# Display help
alchemist --help

# Show version
alchemist --version

# Create a new migration
alchemist make migration create_users_table

# Create a new model
alchemist make model User

# Create a new handler
alchemist make handler UserHandler

# Run pending migrations
alchemist migrate run
```

## 📖 Usage

### Make Commands

The `make` command provides scaffolding capabilities to generate boilerplate code.

#### Creating Migrations

```bash
# Create a new migration with timestamp
alchemist make migration create_users_table

# Specify custom directory
alchemist make migration create_posts_table --dir=database/migrations
```

This creates two files:
- `{timestamp}_create_users_table.up.sql` - Contains the migration SQL
- `{timestamp}_create_users_table.down.sql` - Contains the rollback SQL

#### Creating Models

```bash
# Create a new model
alchemist make model User

# Specify custom directory
alchemist make model Product --dir=internal/domain/models
```

Generated model structure:
```go
package models

type User struct {
    // Fields here
}
```

#### Creating Handlers

```bash
# Create a new HTTP handler
alchemist make handler UserHandler

# Specify custom directory
alchemist make handler ProductHandler --dir=internal/api/handlers
```

Generated handler structure:
```go
package handlers

import (
    "net/http"
)

// UserHandler handles HTTP requests.
func UserHandler(w http.ResponseWriter, r *http.Request) {
    // TODO: Implement handler logic
}
```

### Migration Commands

The `migrate` command manages database migrations using the [golang-migrate](https://github.com/golang-migrate/migrate) library.

#### Running Migrations

```bash
# Run all pending migrations
alchemist migrate run

# Run migrations with custom configuration
alchemist migrate run --config=config.yaml --dir=database/migrations
```

## ⚙️ Configuration

Alchemist supports configuration through multiple methods:

### Environment Variables

Set configuration using environment variables:

```bash
export ALCHEMIST_DB_URL="postgres://user:password@localhost/dbname?sslmode=disable"
export ALCHEMIST_MIGRATIONS_DIR="internal/database/migrations"
```

### Configuration File

Create a configuration file (supports YAML, JSON, TOML):

```yaml
# config.yaml
database:
  url: "postgres://user:password@localhost/dbname?sslmode=disable"
  
migrations:
  dir: "internal/database/migrations"

make:
  models_dir: "internal/models"
  handlers_dir: "internal/handlers"
```

Use the configuration file:

```bash
alchemist --config=config.yaml migrate run
```

### Default Directories

| Component | Default Directory |
|-----------|------------------|
| Migrations | `internal/database/migrations` |
| Models | `internal/models` |
| Handlers | `internal/handlers` |

## 🏗️ Project Structure

```
alchemist/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   └── cli/
│       ├── main.go          # CLI root command
│       ├── make/            # Scaffolding commands
│       │   ├── root.go
│       │   ├── migration.go
│       │   ├── model.go
│       │   └── handler.go
│       ├── migrate/         # Migration commands
│       │   ├── root.go
│       │   └── migrate.go
│       └── response/        # Output formatting
│           └── response.go
├── .goreleaser.yaml         # Release configuration
├── go.mod                   # Go module definition
└── CHANGELOG.md             # Project changelog
```

## 🛠️ Development

### Prerequisites

- Go 1.24.4 or later
- Git

### Building from Source

```bash
# Clone the repository
git clone https://github.com/go-alchemist/alchemist.git
cd alchemist

# Download dependencies
go mod tidy

# Build the application
go build -o alchemist cmd/main.go

# Run the application
./alchemist --help
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -cover ./...
```

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter (if available)
golangci-lint run
```

## 📋 Supported Databases

Alchemist supports the following databases through golang-migrate:

- PostgreSQL
- MySQL
- SQLite
- SQL Server
- And more...

## 🤝 Contributing

We welcome contributions! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Process

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow standard Go conventions
- Use `gofmt` to format your code
- Add comments for exported functions and types
- Write tests for new functionality

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [golang-migrate](https://github.com/golang-migrate/migrate) - Database migrations
- [Fatih Color](https://github.com/fatih/color) - Colored terminal output

## 📞 Support

If you have any questions or run into issues, please:

1. Check the [documentation](README.md)
2. Search [existing issues](https://github.com/go-alchemist/alchemist/issues)
3. Create a [new issue](https://github.com/go-alchemist/alchemist/issues/new)

---

Made with ❤️ by the Alchemist team