package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/SevereCloud/vksdk/v2/api/params"
	iris "github.com/sh2nk/shreks-callback/iris-callback-api"
)

func OnAddUser(ctx context.Context, w http.ResponseWriter, s iris.IrisSignal) {
	b := params.NewMessagesAddChatUserBuilder()
	cp, _ := getChatPair(ctx, s.UserID, s.Object.(iris.AddUser).Chat)

	b.ChatID(cp.ChatID)
	b.UserID(s.Object.(iris.AddUser).UserID)
	b.Params["visible_messages_count"] = 250

	_, err := VK.MessagesAddChatUser(b.Params)
	if err != nil {
		log.Printf("Couldn't add user to chat: %v", err)
		sendMessage(fmt.Sprint(iris.Icons.Warn, "Ошибка добавления: ", err.Error()), cp.ChatID)
	}
}

func OnSubscribeSignals(ctx context.Context, w http.ResponseWriter, s iris.IrisSignal) {
	cp, ok := getChatPair(ctx, s.UserID, s.Object.(iris.SubscribeSignals).Chat)
	if !ok {

	}

	sendMessage(fmt.Sprint(iris.Icons.Success, "Беседа распознана"), cp.ChatID)
}

func sendMessage(m string, chatID int) {
	b := params.NewMessagesSendBuilder()
	b.PeerID(2000000000 + chatID)
	b.RandomID(int(randomInt32()))
	b.Message(m)

	VK.MessagesSend(b.Params)
}
