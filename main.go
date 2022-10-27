package main

import (
	"log"
	"net/http"
	"os"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	iris "github.com/sh2nk/shreks-callback/iris-callback-api"
)

var (
	VK      *api.VK
	VKToken string
	Addr    string
)

// TODO: make proper event system in iris api package.
func callback(w http.ResponseWriter, r *http.Request) {
	signal, err := iris.UnmarshalSignal(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	switch signal.Method {
	case "addUser":
		b := params.NewMessagesAddChatUserBuilder()
		b.ChatID(signal.Object.(iris.AddUser).Chat)
		b.UserID(signal.Object.(iris.AddUser).UserID)
		_, err = VK.MessagesAddChatUser(b.Params)
		if err != nil {
			log.Printf("Couldn't add user to chat: %v", err)
		}
	}
}

// Gets some values from env vars, otherwise returns fallback value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func init() {
	VKToken = getEnv("SHREK_VK_TOKEN", "fallbacktoken")
}

func main() {
	VK = api.NewVK(VKToken)

	http.HandleFunc("/callback", callback)
	log.Printf("Started callback route on %s...", Addr)
	go http.ListenAndServe(Addr, nil)
}
