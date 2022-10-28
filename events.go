package main

import (
	"context"
	"log"
	"net/http"

	"github.com/SevereCloud/vksdk/v2/api/params"
	iris "github.com/sh2nk/shreks-callback/iris-callback-api"
)

func OnPing(ctx context.Context, w http.ResponseWriter, s iris.IrisSignal) {

}

func OnAddUser(ctx context.Context, w http.ResponseWriter, s iris.IrisSignal) {
	b := params.NewMessagesAddChatUserBuilder()
	cp, _ := getChatPair(ctx, s.UserID, s.Object.(iris.AddUser).Chat)

	b.ChatID(cp.ChatID)
	b.UserID(s.Object.(iris.AddUser).UserID)
	b.Params["visible_messages_count"] = 250

	_, err := VK.MessagesAddChatUser(b.Params)
	if err != nil {
		log.Printf("Couldn't add user to chat: %v", err)
	}
}

func OnSubscribeSignals(ctx context.Context, w http.ResponseWriter, s iris.IrisSignal) {
	cp, ok := getChatPair(ctx, s.UserID, s.Object.(iris.SubscribeSignals).Chat)
	if !ok {

	}

	b := params.NewMessagesSendBuilder()
	b.PeerID(2000000000 + cp.ChatID)
	b.RandomID(int(randomInt32()))
	b.Message("✅ Беседа распознана")
}
