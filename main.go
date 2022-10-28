package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/SevereCloud/vksdk/v2/api"
	_ "github.com/lib/pq"
)

var (
	VK      *api.VK
	VKToken string
	Addr    string
	PQURL   string
	DB      *sql.DB
)

func init() {
	VKToken = getEnv("SHREK_VK_TOKEN", "fallbacktoken")
	Addr = getEnv("SHREK_PORT", ":5000")
	PQURL = getEnv("POSTGRES_URL", "postgres://user:password@localhost:5432/s3cr3t")
}

func main() {
	var err error
	VK = api.NewVK(VKToken)

	ctx := context.Background()

	// Init db conection
	DB, err = sql.Open("postgres", PQURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer DB.Close()

	createTables(ctx)

	mux := http.NewServeMux()
	mux.Handle("/callback")

	log.Printf("Started callback route on %s...", Addr)
	http.ListenAndServe(Addr, nil)
}
