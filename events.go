package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SevereCloud/vksdk/v2/api/params"
	iris "github.com/sh2nk/shreks-callback/iris-callback-api"
)

func OnPing(w http.ResponseWriter, s iris.IrisSignal) {
	fmt.Fprint(w, "ok")
}

func OnAddUser(w http.ResponseWriter, s iris.IrisSignal) {
	b := params.NewMessagesAddChatUserBuilder()
	//b.ChatID(s.Object.(iris.AddUser).Chat)
	b.UserID(s.Object.(iris.AddUser).UserID)
	_, err := VK.MessagesAddChatUser(b.Params)
	if err != nil {
		log.Printf("Couldn't add user to chat: %v", err)
	}
}

func OnSubscribeSignals(w http.ResponseWriter, s iris.IrisSignal) {
	fmt.Fprint(w, "ok")
}
