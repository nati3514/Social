# Social Media API

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/nati3514/Social?style=social)](https://github.com/nati3514/Social/stargazers)
[![GitHub Issues](https://img.shields.io/github/issues/nati3514/Social)](https://github.com/nati3514/Social/issues)
[![Go Report Card](https://goreportcard.com/badge/github.com/nati3514/Social)](https://goreportcard.com/report/github.com/nati3514/Social)

A modern, high-performance social media API built with Go, featuring real-time capabilities and a RESTful architecture.

## 🚀 Features

### Current Features (v0.3.0)
- ✅ **Core Infrastructure**
  - Health Check Endpoint
  - Chi Router with middleware stack
  - Environment-based configuration
  - Structured logging
  - Database migrations
  - **Input Validation** with `go-playground/validator`

- ✅ **API Features**
  - **JSON Response Formatting** - Consistent API responses
  - **User Feed** - View posts from followed users
  - **Post Management** - Create and retrieve user posts
  - **Error Handling** - Structured error responses
  - **Request Validation** - Input validation middleware

- 🚧 **In Progress**
  - User authentication
  - Comments system
  - Real-time notifications
- ✅ **Repository Pattern** - Clean data access layer

### Planned Features
- [ ] User authentication & authorization (JWT)
- [ ] User registration & login
- [ ] Post creation, editing, and deletion
- [ ] Comments system
- [ ] Like/Unlike functionality
- [ ] Follow/Unfollow users
- [ ] User profiles
- [ ] Feed generation
- [ ] Search functionality
- [ ] Rate limiting
- [ ] Redis caching
- [ ] File upload (images/videos)
- [ ] Real-time notifications
- [ ] Docker support

## 📋 Prerequisites

- **Go** 1.21 or higher
- **Git** for version control
- **PostgreSQL** 13+
- **Air** (optional, for live reload during development)
- **Golang Migrate** (for database migrations)

## 🛠️ Installation

### 1. Clone the repository
```bash
git clone https://github.com/nati3514/Social.git
cd Social
```

### 2. Install dependencies
```bash
go mod download
```

### 3. Set up environment variables
Create a `.env` file in the root directory:
```bash
# Server
ADDR=:4000

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=social
DB_SSLMODE=disable

# Connection Pool
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
DB_MAX_IDLE_TIME=15m

### 4. Run the application

**Option A: Using Go directly**
```bash
go run ./cmd/api
```

**Option B: Using Air (with live reload)**
```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with Air
air
```

The API will start on the port specified in your `.env` file (default: `:4000`).

## 📁 Project Structure

```
Social/
├── cmd/
│   ├── api/                 # API server
│   │   ├── main.go          # Application entry point
│   │   ├── api.go           # Server setup and routing
│   │   └── health.go        # Health check handler
│   └── migrate/             # Database migrations
│       └── migrations/      # Migration files
│           ├── *.up.sql     # SQL for applying migrations
│           └── *.down.sql   # SQL for rolling back migrations
├── internal/
│   ├── db/                  # Database connection and setup
│   ├── env/                 # Environment variable helpers
│   └── store/               # Repository pattern implementation
│       ├── postgres/        # PostgreSQL implementations
│       └── store.go         # Store interfaces
├── bin/                     # Compiled binaries (gitignored)
├── .env                     # Environment variables (gitignored)
├── .air.toml                # Air configuration
├── migrate.ps1              # Windows migration helper
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
└── README.md                # This file
```

## 🔌 API Endpoints

### Current Endpoints

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| `GET` | `/v1/health` | Health check | ✅ Implemented |

### Example Usage

**Health Check:**
```bash
curl http://localhost:4000/v1/health
```

Response:
```
ok
```

## 🛠 Database Migrations

### Windows
```powershell
# Create a new migration
.\migrate.ps1 create migration_name

# Apply all migrations
.\migrate.ps1 up all

# Rollback last migration
.\migrate.ps1 down 1

# Check migration status
.\migrate.ps1 version
```

### Linux/macOS
```bash
# Install migrate
brew install golang-migrate

# Create a new migration
migrate create -ext sql -dir cmd/migrate/migrations -seq migration_name

# Apply all migrations
migrate -path=cmd/migrate/migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up

# Rollback last migration
migrate -path=cmd/migrate/migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" down 1
```

## 🧪 Development

### Running with Air (Live Reload)
Air automatically reloads the application when you make changes to `.go` files:

```bash
air
```

Configuration is in `.air.toml`.

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ADDR` | Server address and port | `:8080` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | - |
| `DB_NAME` | Database name | `social` |
| `DB_SSLMODE` | SSL mode for database | `disable` |
| `DB_MAX_OPEN_CONNS` | Max open connections | `25` |
| `DB_MAX_IDLE_CONNS` | Max idle connections | `25` |
| `DB_MAX_IDLE_TIME` | Max connection idle time | `15m` |

### Code Style

This project follows [Conventional Commits](https://www.conventionalcommits.org/) for commit messages:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `test:` - Test additions/changes
- `chore:` - Maintenance tasks

**Example commits:**
```bash
feat(auth): add JWT authentication
fix(api): resolve health check timeout
docs: update API endpoint documentation
```

## 🏗️ Built With

- **[Go](https://golang.org/)** - Programming language
- **[Chi](https://github.com/go-chi/chi)** - HTTP router
- **[godotenv](https://github.com/joho/godotenv)** - Environment variable management
- **[Air](https://github.com/cosmtrek/air)** - Live reload for development
- **[validator](https://github.com/go-playground/validator)** - Request payload validation

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Natnael** - [@nati3514](https://github.com/nati3514)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch (`git checkout -b feat/amazing-feature`)
3. Commit your changes using conventional commits (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feat/amazing-feature`)
5. Open a Pull Request

## 📊 Roadmap

### Phase 1: Core API (Current)
- [x] Project setup
- [x] Health check endpoint
- [x] Environment configuration
- [ ] Database setup (PostgreSQL)
- [ ] User authentication

### Phase 2: User Management
- [ ] User registration
- [ ] User login
- [ ] User profiles
- [ ] Password reset

### Phase 3: Social Features
- [ ] Posts (CRUD)
- [ ] Comments
- [ ] Likes
- [ ] Follow system

### Phase 4: Advanced Features
- [ ] Feed algorithm
- [ ] Search
- [ ] Notifications
- [ ] File uploads

### Phase 5: Performance & Scale
- [ ] Caching (Redis)
- [ ] Rate limiting
- [ ] API documentation (Swagger)
- [ ] Docker containerization

## 📞 Support

For questions or issues, please open an issue on GitHub.

---

**⭐ Star this repository if you find it helpful!**
