package main

import (
	"log"

	"github.com/nati3514/Social/internal/db"
	"github.com/nati3514/Social/internal/env"
	"github.com/nati3514/Social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://postgres:12345@localhost:5432/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}
