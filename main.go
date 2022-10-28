package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	_ "github.com/lib/pq"
)

var (
	VK      *api.VK
	VKToken string
	Addr    string
	PQURL   string
	DB      *sql.DB
	UserID  int
)

func init() {
	VKToken = getEnv("SHREK_VK_TOKEN", "fallbacktoken")
	Addr = getEnv("SHREK_PORT", ":5000")
	PQURL = getEnv("POSTGRES_URL", "postgres://user:password@localhost:5432/s3cr3t")
}

func main() {
	var err error
	ctx := context.Background()

	// Connect to VK API
	VK = api.NewVK(VKToken)

	// Get current user ID
	b := params.NewUsersGetBuilder()
	resp, err := VK.UsersGet(b.Params)
	if err != nil {
		log.Fatalf("Unable to get user id: %v\n", err)
	}
	UserID = resp[0].ID

	// Init db conection
	DB, err = sql.Open("postgres", PQURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer DB.Close()

	// DB tables init
	createTables(ctx)

	// Register new userbot
	registerSecret(ctx, UserID, getEnv("USERBOT_SECRET", "s0me-s3r1ous-sh1t"))

	mux := http.NewServeMux()
	mux.Handle("/callback", HandleSignal(Auth(Callback(OK()))))
	log.Printf("Started callback route on %s...", Addr)
	http.ListenAndServe(Addr, mux)
}
