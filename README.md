# ST Bank - Simple Banking API

A RESTful HTTP API for a simple banking system built with Go, featuring user management, accounts, and money transfers with proper transaction handling and security.

## Features

- 🔐 **User Management** - Create users with secure bcrypt password hashing
- 🔑 **Authentication** - JWT and PASETO token-based authentication
- 💰 **Account Management** - Create, read, update, and delete bank accounts
- 💸 **Money Transfers** - Secure transfers between accounts with transaction support
- 🔒 **Database Transactions** - ACID compliance with deadlock prevention
- ✅ **Input Validation** - Custom validators and comprehensive error handling
- 🧪 **Comprehensive Testing** - Unit tests with mock database and 96%+ coverage
- 🐘 **PostgreSQL** - Production-ready database with migrations
- 📊 **Database Schema** - Clean schema with foreign key constraints

## Tech Stack

- **Backend**: Go 1.24.4
- **Web Framework**: Gin
- **Database**: PostgreSQL 12
- **ORM**: SQLC (type-safe SQL code generation)
- **Migrations**: golang-migrate
- **Authentication**: JWT & PASETO tokens
- **Password Hashing**: bcrypt
- **Testing**: Testify + GoMock
- **Configuration**: Viper
- **Containerization**: Docker

## Project Structure

```
st-bank/
├── api/                    # HTTP API handlers and routes
│   ├── account.go         # Account CRUD operations
│   ├── transfer.go        # Money transfer operations  
│   ├── user.go           # User management & authentication
│   ├── validator.go      # Custom validation logic
│   └── *_test.go         # API tests
├── db/
│   ├── migration/        # Database migration files
│   ├── query/           # SQL queries
│   ├── sqlc/           # Generated type-safe Go code
│   └── mock/           # Mock database for testing
├── token/              # Token authentication
│   ├── jwt_maker.go   # JWT token implementation
│   ├── jwt_maker_test.go # JWT comprehensive tests
│   ├── paseto_maker.go # PASETO token implementation
│   ├── paseto_maker_test.go # PASETO comprehensive tests
│   ├── payload.go     # Token payload structure
│   └── maker.go       # Token interface
├── util/               # Utility functions
│   ├── config.go      # Configuration management
│   ├── password.go    # Password hashing utilities
│   └── random.go      # Test data generation
├── main.go            # Application entry point
├── Makefile          # Build and development tasks
└── app.env           # Environment configuration
```

## Database Schema

### Users Table
- `username` (PK) - Unique username
- `hashed_password` - bcrypt hashed password
- `full_name` - User's full name
- `email` - Unique email address
- `password_changed_at` - Timestamp of last password change
- `created_at` - Account creation timestamp

### Accounts Table
- `id` (PK) - Account ID
- `owner` (FK) - References users.username
- `balance` - Account balance in cents
- `currency` - Currency code (USD, EUR, UAH)
- `created_at` - Account creation timestamp
- **Unique constraint**: (owner, currency) - One account per currency per user

### Transfers Table
- `id` (PK) - Transfer ID
- `from_account_id` (FK) - Source account
- `to_account_id` (FK) - Destination account
- `amount` - Transfer amount (must be positive)
- `created_at` - Transfer timestamp

### Entries Table
- `id` (PK) - Entry ID
- `account_id` (FK) - Related account
- `amount` - Entry amount (can be negative or positive)
- `created_at` - Entry timestamp

## API Endpoints

### Authentication
- `POST /users` - Create a new user
- `POST /users/login` - Login user and get access token

### Accounts
- `POST /accounts` - Create a new account
- `GET /accounts/:id` - Get account by ID
- `GET /accounts` - List accounts (paginated)
- `PUT /accounts/:id` - Update account balance
- `DELETE /accounts/:id` - Delete account

### Transfers
- `POST /transfers` - Transfer money between accounts

## Getting Started

### Prerequisites

- Go 1.24.4+
- PostgreSQL 12+
- Docker (optional)
- golang-migrate CLI tool

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/volskyi-dmytro/st-bank.git
   cd st-bank
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Install golang-migrate**
   ```bash
   # On macOS
   brew install golang-migrate

   # On Linux/Windows - download from releases
   # https://github.com/golang-migrate/migrate/releases
   ```

4. **Set up PostgreSQL**
   ```bash
   # Using Docker
   make postgres
   make createdb

   # Or use your own PostgreSQL instance
   createdb st_bank
   ```

5. **Run database migrations**
   ```bash
   make migrateup
   ```

6. **Configure environment**
   ```bash
   cp app.env.example app.env
   # Edit app.env with your database credentials
   ```

### Running the Application

1. **Start the server**
   ```bash
   make server
   # or
   go run main.go
   ```

2. **Run tests**
   ```bash
   make test
   ```

3. **Generate SQLC code** (if you modify SQL queries)
   ```bash
   make sqlc
   ```

## Usage Examples

### Create a User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "secret123",
    "full_name": "Alice Johnson", 
    "email": "alice@example.com"
  }'
```

### Login User
```bash
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "secret123"
  }'
```

Response:
```json
{
  "access_token": "v2.local.Gdh5kiOTyyaQ3_bNykYDeYHO21Jg2...",
  "user": {
    "username": "alice",
    "full_name": "Alice Johnson",
    "email": "alice@example.com",
    "password_changed_at": "2023-01-01T00:00:00Z",
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### Create an Account
```bash
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{
    "owner": "alice",
    "currency": "USD"
  }'
```

### Transfer Money
```bash
curl -X POST http://localhost:8080/transfers \
  -H "Content-Type: application/json" \
  -d '{
    "from_account_id": 1,
    "to_account_id": 2,
    "amount": 1000,
    "currency": "USD"
  }'
```

### Get Account
```bash
curl http://localhost:8080/accounts/1
```

## Configuration

The application uses environment variables for configuration. See `app.env`:

```env
DB_DRIVER=postgres
DB_SOURCE=postgresql://root:password@localhost:5432/st_bank?sslmode=disable
SERVER_ADDRESS=0.0.0.0:8080
TOKEN_SYMMETRIC_KEY=12345678901234567890123456789012
ACCESS_TOKEN_DURATION=15m
TOKEN_TYPE=paseto
```

### Token Configuration

- `TOKEN_SYMMETRIC_KEY`: 32-character symmetric key for token signing/encryption
- `ACCESS_TOKEN_DURATION`: Token expiration time (e.g., 15m, 1h, 24h)
- `TOKEN_TYPE`: Choose between `jwt` or `paseto` (default: paseto)

#### JWT vs PASETO

- **JWT**: JSON Web Tokens - widely adopted, good tooling support
- **PASETO**: Platform-Agnostic Security Tokens - more secure by design, prevents common JWT vulnerabilities

Switch between token types by changing `TOKEN_TYPE` in `app.env`.

## Development

### Available Make Commands

```bash
make postgres        # Start PostgreSQL container
make createdb        # Create database
make dropdb         # Drop database
make migrateup      # Run all migrations
make migratedown    # Rollback all migrations
make migrateup1     # Run one migration
make migratedown1   # Rollback one migration
make sqlc           # Generate SQLC code
make test           # Run all tests
make server         # Start the server
```

### Testing

The project includes comprehensive tests:

- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test API endpoints with mock database
- **Database Tests**: Test SQL queries against real database
- **Coverage**: 96%+ test coverage

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -cover ./...

# Run specific test
go test -v ./api -run TestTransferAPI

# Run token authentication tests
go test -v ./token -run TestJWT
go test -v ./token -run TestPaseto
```

### Token Authentication Testing

The project includes comprehensive test coverage for both JWT and PASETO token implementations:

#### JWT Tests (`jwt_maker_test.go`)
- ✅ **Core Functionality**: Token creation, verification, and payload validation
- ✅ **Security Tests**: Algorithm confusion prevention, signature validation
- ✅ **Error Handling**: Expired tokens, malformed tokens, invalid keys
- ✅ **Edge Cases**: Empty usernames, zero/long durations, key size validation
- ✅ **Attack Prevention**: None algorithm bypass, cross-key verification failure

#### PASETO Tests (`paseto_maker_test.go`)
- ✅ **Core Functionality**: Encryption, decryption, and format validation
- ✅ **Security Tests**: Key isolation, token uniqueness, ChaCha20-Poly1305 validation
- ✅ **Error Handling**: Expired tokens, malformed tokens, wrong key sizes
- ✅ **Edge Cases**: Various key sizes, token format validation
- ✅ **Stress Testing**: 1000+ token creation/verification cycles
- ✅ **Advanced Security**: Nonce randomness, multi-key isolation

Both test suites ensure resistance against common token-based attacks and validate proper implementation of each token standard's security features.

```bash
# Run all token tests
go test -v ./token

# Test JWT implementation only
go test -v ./token -run TestJWT

# Test PASETO implementation only  
go test -v ./token -run TestPaseto
```

## Security Features

- ✅ **Password Hashing**: bcrypt with salt
- ✅ **Token Authentication**: JWT and PASETO support with configurable expiration
- ✅ **Input Validation**: Comprehensive request validation
- ✅ **SQL Injection Prevention**: Parameterized queries via SQLC
- ✅ **Transaction Safety**: ACID compliance with deadlock prevention
- ✅ **Error Handling**: Secure error responses without information leakage

## Database Migrations

Migrations are managed with golang-migrate:

```bash
# Create new migration
migrate create -ext sql -dir db/migration -seq add_users

# Apply migrations
make migrateup

# Rollback migrations  
make migratedown
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Testing in CI/CD

The project includes GitHub Actions workflow that:
- Sets up PostgreSQL service
- Runs database migrations
- Executes all tests
- Validates code coverage

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built following clean architecture principles
- Inspired by modern banking system requirements
- Uses industry-standard security practices

---

**Note**: This is a learning project for demonstrating Go backend development best practices. Not intended for production banking use without additional security measures.