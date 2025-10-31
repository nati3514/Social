# Social Media API

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/nati3514/Social?style=social)](https://github.com/nati3514/Social/stargazers)
[![GitHub Issues](https://img.shields.io/github/issues/nati3514/Social)](https://github.com/nati3514/Social/issues)
[![Go Report Card](https://goreportcard.com/badge/github.com/nati3514/Social)](https://goreportcard.com/report/github.com/nati3514/Social)

A modern, high-performance social media API built with Go, featuring real-time capabilities and a RESTful architecture.

## 🚀 Features

### Current Features (v0.4.0)
- ✅ **Core Infrastructure**
  - Health Check Endpoint
  - Chi Router with context middleware
  - Environment-based configuration
  - Structured logging with request IDs
  - Database migrations with versioning
  - **Multi-layer Validation** with `go-playground/validator`

- ✅ **API Features**
  - **RESTful JSON API** with consistent response format
  - **Post Management** - Full CRUD operations
  - **User Profiles** - View user details and activity
  - **Comment System** - Nested comments on posts
  - **Context Middleware** - Efficient resource loading
  - **Error Handling** - Structured error responses with proper HTTP status codes
  - **Request Validation** - Input validation at multiple layers
  - **Partial Updates** - PATCH support with proper null handling
  - **Optimistic Concurrency Control** - Version-based updates to prevent lost updates

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
│   └── migrate/             # Database migration and seeding
│       └── seed/            # Database seeding
│           └── main.go      # Seed data generation
├── migrations/              # Database migration files
│   ├── *.up.sql            # SQL for applying migrations
│   └── *.down.sql          # SQL for rolling back migrations
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

### API Endpoints

#### Posts
| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| `GET` | `/v1/posts` | List all posts | ✅ Implemented |
| `POST` | `/v1/posts` | Create a new post | ✅ Implemented |
| `GET` | `/v1/posts/{id}` | Get a specific post with comments | ✅ Implemented |
| `PATCH` | `/v1/posts/{id}` | Partially update a post | ✅ Implemented |
| `DELETE` | `/v1/posts/{id}` | Delete a post | ✅ Implemented |

#### Comments
| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| `GET` | `/v1/posts/{id}/comments` | Get comments for a post | ✅ Implemented |

#### System
| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| `GET` | `/v1/health` | Health check | ✅ Implemented |

#### Users
| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| `GET` | `/v1/users/{userID}` | Get user profile | ✅ Implemented |
| `PUT` | `/v1/users/{userID}/follow` | Follow a user | ✅ Implemented |
| `PUT` | `/v1/users/{userID}/unfollow` | Unfollow a user | ✅ Implemented |

### Example Usage

**Create a Post:**
```bash
curl -X POST http://localhost:4000/v1/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"First Post","content":"This is my first post","user_id":1}'
```

**Get a Post with Comments:**
```bash
curl http://localhost:4000/v1/posts/1
```

**Update a Post (Partial Update):**
```bash
curl -X PATCH http://localhost:4000/v1/posts/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Title"}'
```

**Update a Post with Version Control:**
```bash
curl -X PATCH http://localhost:4000/v1/posts/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Title", "version": 5}'
```

**Health Check:**
```bash
curl http://localhost:4000/v1/health
```

**User Operations:**

```bash
# Get user by ID
curl http://localhost:4000/v1/users/1

# Follow a user
curl -X PUT http://localhost:4000/v1/users/2/follow \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1}'

# Unfollow a user
curl -X PUT http://localhost:4000/v1/users/2/unfollow \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1}'

# Successful response (204 No Content for follow/unfollow)
# No content in response body

# Error responses:
# Already following (409 Conflict)
{
  "error": "record already exists"
}

# User not found (404 Not Found)
{
  "error": "record not found"
}

# Get user profile response (200 OK)
{
  "data": {
    "id": 1,
    "name": "Nati Age",
    "email": "nati@example.com",
    "created_at": "2023-10-31T10:00:00Z"
  }
}
```

## 🔄 Concurrency Control

The API implements optimistic concurrency control to handle concurrent updates to resources. When updating a post, include the current version number in the request. If the version on the server doesn't match the provided version, the update will be rejected with a `409 Conflict` status code.

### How It Works
1. Each post has a version number that increments with each update
2. When updating a post, include the current version in the request
3. The server verifies the version matches before applying updates
4. If versions don't match, a 409 Conflict is returned

### Error Response (409 Conflict)
```json
{
    "error": "edit conflict: post has been modified by another user"
}
```

## Database Seeding

The application includes a database seeding system to generate test data for development and testing. The seeder creates realistic-looking users, posts, and comments with varied content.

### Running the Seeder

#### Using PowerShell:
```powershell
.\migrate.ps1 seed
```

#### Using Make:
```bash
make seed
```

### Seeding Details
- Creates 100 random users with realistic names and email addresses
- Generates 20 posts with varied content and tags
- Creates 20 comments on random posts
- Includes proper error handling and logging

### Example Output
```bash
Seeding database...
2025/10/30 17:12:57 Database seeded successfully
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

# Seed the database with test data
.\migrate.ps1 seed
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

# Seed the database with test data
make seed
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

### Code Style & Best Practices

#### Commit Messages
This project follows [Conventional Commits](https://www.conventionalcommits.org/) for commit messages:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `test:` - Test additions/changes
- `chore:` - Maintenance tasks

#### Code Organization
- **Handler Layer**: HTTP request/response handling, input validation
- **Service Layer**: Business logic, data validation
- **Store Layer**: Database operations, data integrity
- **Middleware**: Cross-cutting concerns (logging, auth, context management)

#### Error Handling
- Consistent error responses with appropriate HTTP status codes
- Detailed error messages in development
- Secure error messages in production
- Structured error logging

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
