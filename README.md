# ST Bank - Simple Banking API

A RESTful HTTP API for a simple banking system built with Go, featuring user management, accounts, and money transfers with proper transaction handling and security.

## Features

- 🔐 **User Management** - Create users with secure bcrypt password hashing
- 🔑 **Authentication** - JWT and PASETO token-based authentication with middleware
- 🛡️ **Authorization** - Role-based access control and ownership validation
- 💰 **Account Management** - Create, read, update, and delete bank accounts
- 💸 **Money Transfers** - Secure transfers between accounts with transaction support
- 🔒 **Database Transactions** - ACID compliance with deadlock prevention
- ✅ **Input Validation** - Custom validators and comprehensive error handling
- 🧪 **Comprehensive Testing** - Unit tests with mock database and 96%+ coverage
- 🐘 **PostgreSQL** - Production-ready database with migrations
- 📊 **Database Schema** - Clean schema with foreign key constraints

## Tech Stack

- **Backend**: Go 1.24.5
- **Web Framework**: Gin
- **Database**: PostgreSQL 12
- **ORM**: SQLC (type-safe SQL code generation)
- **Migrations**: golang-migrate
- **Authentication**: JWT & PASETO tokens
- **Password Hashing**: bcrypt
- **Testing**: Testify + GoMock
- **Configuration**: Viper
- **Containerization**: Docker & Docker Compose

## Project Structure

```
st-bank/
├── api/                    # HTTP API handlers and routes
│   ├── account.go         # Account CRUD operations
│   ├── transfer.go        # Money transfer operations  
│   ├── user.go           # User management & authentication
│   ├── middleware.go      # Authentication middleware
│   ├── validator.go      # Custom validation logic
│   └── *_test.go         # API tests
├── db/
│   ├── migration/        # Database migration files
│   ├── query/           # SQL queries
│   ├── sqlc/           # Generated type-safe Go code
│   └── mock/           # Mock database for testing
├── eks/                  # Kubernetes deployment manifests
│   ├── aws-auth.yaml    # EKS authentication ConfigMap
│   ├── deployment.yaml  # Kubernetes Deployment
│   └── service.yaml     # LoadBalancer Service
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
├── Dockerfile        # Docker container configuration
├── docker-compose.yaml # Multi-container orchestration
├── start.sh          # Container startup script
├── wait-for.sh       # Database readiness checker
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

### Authentication (Public)
- `POST /users` - Create a new user
- `POST /users/login` - Login user and get access token

### Accounts (Protected) 🔒
- `POST /accounts` - Create a new account (requires authentication)
- `GET /accounts/:id` - Get account by ID (requires authentication + ownership)
- `GET /accounts` - List accounts (requires authentication, filtered by owner)
- `PUT /accounts/:id` - Update account balance (requires authentication + ownership)
- `DELETE /accounts/:id` - Delete account (requires authentication + ownership)

### Transfers (Protected) 🔒
- `POST /transfers` - Transfer money between accounts (requires authentication + ownership of source account)

**Authentication Required**: All protected endpoints require a valid Bearer token in the Authorization header.
**Authorization**: Users can only access and modify their own accounts and transfers.

## Getting Started

### Prerequisites

- Go 1.24.5+
- PostgreSQL 12+ (or Docker)
- Docker & Docker Compose (recommended)
- golang-migrate CLI tool (if not using Docker)

### Installation

#### Option 1: Docker Compose (Recommended) 🐳

1. **Clone the repository**
   ```bash
   git clone https://github.com/volskyi-dmytro/st-bank.git
   cd st-bank
   ```

2. **Start with Docker Compose**
   ```bash
   docker compose up --build
   ```

   This will:
   - Build the application image
   - Start PostgreSQL database
   - Wait for database to be ready
   - Run database migrations automatically
   - Start the API server on port 8080

3. **Stop the services**
   ```bash
   docker compose down
   ```

#### Option 2: Local Development

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

#### With Docker Compose
```bash
# Start all services
docker compose up

# Start in background
docker compose up -d

# Rebuild and start
docker compose up --build

# View logs
docker compose logs -f api

# Stop services
docker compose down
```

#### Local Development
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
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "owner": "alice",
    "currency": "USD"
  }'
```

### Transfer Money
```bash
curl -X POST http://localhost:8080/transfers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "from_account_id": 1,
    "to_account_id": 2,
    "amount": 1000,
    "currency": "USD"
  }'
```

### Get Account
```bash
curl -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  http://localhost:8080/accounts/1
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

### Docker Environment

When using Docker Compose, environment variables are automatically configured:
- Database host changes from `localhost` to `postgresdb` (container service name)
- All required environment variables are set in `docker-compose.yaml`
- Database data persists in Docker volume `postgres_data`

### Token Configuration

- `TOKEN_SYMMETRIC_KEY`: 32-character symmetric key for token signing/encryption
- `ACCESS_TOKEN_DURATION`: Token expiration time (e.g., 15m, 1h, 24h)
- `TOKEN_TYPE`: Choose between `jwt` or `paseto` (default: paseto)

#### JWT vs PASETO

- **JWT**: JSON Web Tokens - widely adopted, good tooling support
- **PASETO**: Platform-Agnostic Security Tokens - more secure by design, prevents common JWT vulnerabilities

Switch between token types by changing `TOKEN_TYPE` in `app.env`.

## Development

### Available Commands

#### Docker Commands
```bash
docker compose up --build    # Build and start all services
docker compose up -d         # Start services in background
docker compose down          # Stop and remove containers
docker compose logs -f api   # Follow API logs
docker compose exec api sh   # Access API container shell
docker compose exec postgresdb psql -U root -d st_bank  # Access database
```

#### Make Commands (Local Development)
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
- **Authentication Tests**: Test middleware and token validation
- **Authorization Tests**: Test ownership-based access control
- **Database Tests**: Test SQL queries against real database
- **Coverage**: 96%+ test coverage

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -cover ./...

# Run specific test
go test -v ./api -run TestTransferAPI

# Run authentication tests
go test -v ./api -run TestGetAccountAPI
go test -v ./api -run TestCreateAccountAPI

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
- ✅ **Authentication Middleware**: Automatic token validation for protected routes
- ✅ **Authorization**: Ownership-based access control for all account operations
- ✅ **Route Protection**: Public and protected endpoint separation
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

## Docker Details

### Multi-stage Dockerfile
- **Build Stage**: Compiles Go application and downloads migration tool
- **Runtime Stage**: Minimal Alpine Linux image with only required dependencies
- **Size**: Final image ~56MB
- **Security**: Non-root user, minimal attack surface

### Container Features
- 🔄 **Automatic Migrations**: Database schema applied on startup
- ⏳ **Health Checks**: Wait for database readiness before starting API
- 📁 **Volume Persistence**: Database data survives container restarts
- 🔧 **Environment Flexibility**: Easy configuration via environment variables

### Services
- **API**: Go application on port 8080
- **PostgreSQL**: Database on port 5432 with persistent volume
- **Automatic Migration**: Runs on container startup

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Workflow
```bash
# Start development environment
docker compose up -d

# Make changes to code
# ...

# Rebuild and restart API only
docker compose up --build api

# Run tests
make test

# Clean up
docker compose down
```

## Deployment

### AWS Infrastructure

The application is deployed on AWS with the following components:

- **Database**: Amazon RDS PostgreSQL instance in `us-east-1`
- **Secrets Management**: AWS Secrets Manager for environment variables and database credentials
- **Container Registry**: Amazon ECR in `eu-central-1` for Docker images
- **Container Orchestration**: Amazon EKS (Elastic Kubernetes Service) for production deployment
- **CI/CD**: GitHub Actions automated deployment pipeline

### Deployment Pipeline

The project includes GitHub Actions workflow (`.github/workflows/deploy.yml`) that:

1. **Secrets Loading**: 
   - Connects to AWS Secrets Manager in `us-east-1` region
   - Loads database credentials and application configuration
   - Creates `app.env` file with proper environment variables

2. **Container Build**:
   - Builds Docker image with loaded secrets
   - Tags image with Git commit SHA
   - Pushes to Amazon ECR in `eu-central-1` region

3. **Error Handling**:
   - Validates that secrets are properly loaded (non-empty app.env)
   - Fails deployment if secrets loading fails
   - Provides clear error messages for debugging

### AWS Configuration

**Secrets Manager Configuration:**
```json
{
  "DB_SOURCE": "postgresql://root:PASSWORD@stbank.cw1iu28ck5cv.us-east-1.rds.amazonaws.com:5432/st_bank",
  "DB_DRIVER": "postgres", 
  "SERVER_ADDRESS": "0.0.0.0:8080",
  "ACCESS_TOKEN_DURATION": "15m",
  "TOKEN_SYMMETRIC_KEY": "your-32-character-key"
}
```

**Required AWS Permissions:**
- `secretsmanager:GetSecretValue` for loading application secrets
- `ecr:GetAuthorizationToken`, `ecr:BatchCheckLayerAvailability`, `ecr:GetDownloadUrlForLayer`, `ecr:BatchGetImage` for ECR access

### Kubernetes Deployment

The application includes Kubernetes manifests in the `eks/` directory for EKS deployment:

- **`deployment.yaml`**: Kubernetes Deployment with 2 replicas for high availability
- **`service.yaml`**: LoadBalancer service exposing the API on port 80
- **`aws-auth.yaml`**: ConfigMap for EKS node group and GitHub CI user authentication

```bash
# Deploy to EKS cluster
kubectl apply -f eks/aws-auth.yaml
kubectl apply -f eks/deployment.yaml
kubectl apply -f eks/service.yaml

# Check deployment status
kubectl get pods -l app=stbank-api
kubectl get service stbank-api-service
```

### Running Deployed Image

```bash
# Login to ECR
aws ecr get-login-password --region eu-central-1 | docker login --username AWS --password-stdin 986629373383.dkr.ecr.eu-central-1.amazonaws.com

# Pull and run the latest image
docker pull 986629373383.dkr.ecr.eu-central-1.amazonaws.com/stbank:COMMIT_SHA
docker run -p 8080:8080 986629373383.dkr.ecr.eu-central-1.amazonaws.com/stbank:COMMIT_SHA
```

The deployed application automatically:
- Loads environment variables from the baked-in app.env file
- Runs database migrations on startup
- Connects to the AWS RDS PostgreSQL instance
- Serves the API on port 8080

## Testing in CI/CD

The project includes comprehensive testing and deployment automation:
- Automated secrets loading from AWS Secrets Manager
- Multi-region AWS setup (secrets in us-east-1, ECR in eu-central-1)
- Docker image building and pushing to ECR
- Environment variable validation and error handling

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built following clean architecture principles
- Inspired by modern banking system requirements
- Uses industry-standard security practices

---

**Note**: This is a learning project for demonstrating Go backend development best practices. Not intended for production banking use without additional security measures.
