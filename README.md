# Social Media API

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/nati3514/Social?style=social)](https://github.com/nati3514/Social/stargazers)
[![GitHub Issues](https://img.shields.io/github/issues/nati3514/Social)](https://github.com/nati3514/Social/issues)
[![Go Report Card](https://goreportcard.com/badge/github.com/nati3514/Social)](https://goreportcard.com/report/github.com/nati3514/Social)

A modern, high-performance social media API built with Go, featuring real-time capabilities and a RESTful architecture.

## 🚀 Features

### Current Features (v0.1.0)
- ✅ **Health Check Endpoint** - Monitor API status
- ✅ **Chi Router** - Fast, lightweight HTTP router
- ✅ **Middleware Stack**:
  - Request logging
  - Panic recovery
  - Real IP detection
- ✅ **Environment Configuration** - Flexible config via `.env` files
- ✅ **Live Reload** - Development hot-reload with Air
- ✅ **HTTP Timeouts** - Read/Write timeout protection

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
- [ ] PostgreSQL integration
- [ ] Redis caching
- [ ] File upload (images/videos)
- [ ] Real-time notifications

## 📋 Prerequisites

- **Go** 1.21 or higher
- **Git** for version control
- **Air** (optional, for live reload during development)

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
ADDR=:4000
```

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
│   └── api/
│       ├── main.go          # Application entry point
│       ├── api.go           # Server setup and routing
│       └── health.go        # Health check handler
├── internal/
│   └── env/
│       └── env.go           # Environment variable helpers
├── bin/                     # Compiled binaries (gitignored)
├── .env                     # Environment variables (gitignored)
├── .air.toml               # Air configuration
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
└── README.md               # This file
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
