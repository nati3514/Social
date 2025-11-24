package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/nati3514/Social/internal/db"
	"github.com/nati3514/Social/internal/env"
	"github.com/nati3514/Social/internal/store"
)

const version = "0.0.1"

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host						petstore.swagger.io
// @BasePath					/v2
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	// Initialize configuration
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         "postgres://postgres:12345@localhost:5432/social?sslmode=disable",
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  time.Duration(env.GetInt("DB_MAX_IDLE_MINUTES", 15)) * time.Minute,
		},
		env: env.GetString("ENV", "development"),
	}

	// Initialize database connection
	dbPool, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime.String(),
	)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer func() {
		dbPool.Close()
		log.Println("Database connection closed")
	}()

	log.Println("Database connection pool established")

	// Initialize storage
	storage := store.NewStorage(dbPool)

	// Initialize application
	app := &application{
		config: cfg,
		store:  storage,
	}

	// Setup routes and start server
	router := app.mount()
	log.Printf("Starting server on %s\n", cfg.addr)

	if err := app.run(router); err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}
