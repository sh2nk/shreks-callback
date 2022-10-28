package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/SevereCloud/vksdk/v2/api/params"
	iris "github.com/sh2nk/shreks-callback/iris-callback-api"
)

func OnDeleteMessages(ctx context.Context, w http.ResponseWriter) {
	// Get signal value from conetxt
	s := ctx.Value(signalKey).(iris.IrisSignal)
	o := s.Object.(iris.DeleteMessages)

	b := params.NewMessagesDeleteBuilder()
	cp, _ := getChatPair(ctx, s.UserID, o.Chat)

	b.DeleteForAll(true)
	b.Spam(o.IsSpam)
	b.MessageIDs(o.LocalID)

	_, err := VK.MessagesDelete(b.Params)
	if err != nil {
		log.Printf("Couldn't delete messages: %v", err)
		sendMessage(fmt.Sprint(iris.Icons.Warn, "Не вышло удалить сообщения: ", err.Error()), cp.ChatID)
		return
	}

	sendMessage(fmt.Sprint(iris.Icons.SuccessOff, "Подчистили гавнецо"), cp.ChatID)
}

func OnDeleteMessagesFromUser(ctx context.Context, w http.ResponseWriter) {
	// Get signal value from conetxt
	s := ctx.Value(signalKey).(iris.IrisSignal)
	o := s.Object.(iris.DeleteMessagesFromUser)

	cp, _ := getChatPair(ctx, s.UserID, o.Chat)

	var msg []int

	if o.Amount == 0 {
		o.Amount = 1
	}

	for len(msg) < o.Amount {
		b := params.NewMessagesGetHistoryBuilder()
		b.Count(200)
		b.PeerID(2000000000 + cp.ChatID)

		resp, err := VK.MessagesGetHistory(b.Params)
		if err != nil {
			log.Printf("Couldn't obtain message history: %v", err)
			sendMessage(fmt.Sprint(iris.Icons.Warn, "Не вышло получить историю сообщений: ", err.Error()), cp.ChatID)
			return
		}

		for _, m := range resp.Items {
			if (m.FromID == o.UserID) && (len(msg) < o.Amount) {
				msg = append(msg, m.ID)
			}
		}
	}

	b := params.NewMessagesDeleteBuilder()
	b.DeleteForAll(true)
	b.Spam(o.IsSpam)
	fmt.Println(len(msg))
	b.MessageIDs(msg)

	_, err := VK.MessagesDelete(b.Params)
	if err != nil {
		log.Printf("Couldn't delete messages: %v", err)
		sendMessage(fmt.Sprint(iris.Icons.Warn, "Не вышло удалить сообщения: ", err.Error()), cp.ChatID)
		return
	}

	sendMessage(fmt.Sprint(iris.Icons.SuccessOff, "Подчистили гавнецо"), cp.ChatID)
}

func OnBanExpired(ctx context.Context, w http.ResponseWriter) {
	// Get signal value from conetxt
	s := ctx.Value(signalKey).(iris.IrisSignal)
	o := s.Object.(iris.BanExpired)

	b := params.NewMessagesAddChatUserBuilder()
	cp, _ := getChatPair(ctx, s.UserID, o.Chat)

	b.ChatID(cp.ChatID)
	b.UserID(o.UserID)
	b.Params["visible_messages_count"] = 250

	_, err := VK.MessagesAddChatUser(b.Params)
	if err != nil {
		log.Printf("Couldn't add user to chat: %v", err)
		sendMessage(fmt.Sprint(iris.Icons.Warn, "Ошибка добавления: ", err.Error()), cp.ChatID)
		return
	}

	sendMessage(fmt.Sprint(iris.Icons.Info, "Срок бана пользователя истек."), cp.ChatID)
}

func OnAddUser(ctx context.Context, w http.ResponseWriter) {
	// Get signal value from conetxt
	s := ctx.Value(signalKey).(iris.IrisSignal)
	o := s.Object.(iris.AddUser)

	b := params.NewMessagesAddChatUserBuilder()
	cp, _ := getChatPair(ctx, s.UserID, o.Chat)

	b.ChatID(cp.ChatID)
	b.UserID(o.UserID)
	b.Params["visible_messages_count"] = 250

	_, err := VK.MessagesAddChatUser(b.Params)
	if err != nil {
		log.Printf("Couldn't add user to chat: %v", err)
		sendMessage(fmt.Sprint(iris.Icons.Warn, "Ошибка добавления: ", err.Error()), cp.ChatID)
		return
	}

	sendMessage(fmt.Sprint(iris.Icons.Success, "Вернули пользователя"), cp.ChatID)
}

func OnSubscribeSignals(ctx context.Context, w http.ResponseWriter) {
	// Get signal value from conetxt
	s := ctx.Value(signalKey).(iris.IrisSignal)
	o := s.Object.(iris.SubscribeSignals)

	cp, ok := getChatPair(ctx, s.UserID, o.Chat)
	if !ok {
		chat, err := searchChat(s.Message.Date, s.UserID, o.Text, 10)
		if err != nil {
			log.Printf("Falied to subscribe userbot: %v", err)
			return
		}
		cp = ChatPair{
			UserID:   s.UserID,
			ChatCode: o.Chat,
			ChatID:   chat,
		}
		setChatPair(ctx, cp)
	}

	sendMessage(fmt.Sprint(iris.Icons.Success, "Беседа распознана"), cp.ChatID)
}

func searchChat(date int, userID int, text string, depth int) (int, error) {
	b := params.NewMessagesGetConversationsBuilder()
	b.Count(depth)

	resp, err := VK.MessagesGetConversations(b.Params)
	if err != nil {
		return 0, err
	}
	for _, chat := range resp.Items {
		if chat.Conversation.Peer.Type == "chat" {
			b := params.NewMessagesGetHistoryBuilder()
			b.Count(depth)
			b.PeerID(chat.Conversation.Peer.ID)

			resp, err := VK.MessagesGetHistory(b.Params)
			if err != nil {
				return 0, err
			}

			for _, msg := range resp.Items {
				if (msg.Date == date) && (msg.FromID == userID) && (msg.Text == text) {
					return chat.Conversation.Peer.LocalID, nil
				}
			}
		}
	}
	return 0, errors.New("chat not found")
}

func sendMessage(m string, chatID int) {
	b := params.NewMessagesSendBuilder()
	b.PeerID(2000000000 + chatID)
	b.RandomID(int(randomInt32()))
	b.Message(m)

	VK.MessagesSend(b.Params)
}
