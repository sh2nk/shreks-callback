package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/SevereCloud/vksdk/v2/api"
	_ "github.com/lib/pq"
	iris "github.com/sh2nk/shreks-callback/iris-callback-api"
)

var (
	VK      *api.VK
	VKToken string
	Addr    string
	PQURL   string
	DB      *sql.DB
)

// TODO: make proper event system in iris api package.
func callback(w http.ResponseWriter, r *http.Request) {
	signal, err := iris.UnmarshalSignal(r.Body)
	if err != nil {
		log.Printf("Unmarshal signal error: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	log.Printf("Got Iris %s signal!\n%v\n", signal.Method, signal)

	switch signal.Method {
	case "ping":
		OnPing(w, signal)
	case "addUser":
		OnAddUser(w, signal)
	case "subscribeSignals":
		OnSubscribeSignals(w, signal)
	}
}

func init() {
	VKToken = getEnv("SHREK_VK_TOKEN", "fallbacktoken")
	Addr = getEnv("SHREK_PORT", ":5000")
	PQURL = getEnv("POSTGRES_URL", "postgres://user:password@localhost:5432/s3cr3t")
}

func main() {
	VK = api.NewVK(VKToken)

	ctx := context.Background()

	// Init db conection
	DB, err := sql.Open("postgres", PQURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer DB.Close()

	createTables(ctx)

	http.HandleFunc("/callback", callback)
	log.Printf("Started callback route on %s...", Addr)
	http.ListenAndServe(Addr, nil)
}
