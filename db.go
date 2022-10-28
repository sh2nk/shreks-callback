package main

import (
	"context"
	"database/sql"
	"log"
)

type ChatPair struct {
	UserID   int
	ChatCode string
	ChatID   int
}

func createTables(ctx context.Context) {
	query, err := DB.PrepareContext(ctx, "CREATE TABLE IF NOT EXISTS chatpairs (userID integer primary key, chatCode text, chatID integer);")
	if err != nil {
		log.Fatalf("Could not prepare db create query:%v\n", err)
	}
	_, err = query.ExecContext(ctx)
	if err != nil {
		log.Fatalf("Unable to create database: %v\n", err)
	}
	defer query.Close()

	query, err = DB.PrepareContext(ctx, "CREATE TABLE IF NOT EXISTS userbots (userID integer primary key, secret text);")
	if err != nil {
		log.Fatalf("Could not prepare db create query:%v\n", err)
	}
	_, err = query.ExecContext(ctx)
	if err != nil {
		log.Fatalf("Unable to create database: %v\n", err)
	}
	defer query.Close()
}

func getChatPair(ctx context.Context, userID int, chat string) (ChatPair, bool) {
	query, err := DB.PrepareContext(ctx, "SELECT * FROM chatpairs WHERE userID = $1 AND chatCode = $2;")
	if err != nil {
		log.Fatalf("Could not prepare getChatPair query:%v\n", err)
	}
	res := query.QueryRowContext(ctx, userID, chat)

	var cp ChatPair
	err = res.Scan(&cp.UserID, &cp.ChatCode, &cp.ChatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ChatPair{}, false
		} else {
			log.Fatalf("Chat pair scan error:%v\n", err)
		}
	}
	return cp, true
}

func setChatPair(ctx context.Context, cp ChatPair) {
	query, err := DB.PrepareContext(ctx, `
	INSERT INTO chatpairs (userID, chatCode, chatID)
	VALUES ($1, $2, $3)
	ON CONFLICT (userID) DO UPDATE
		SET chatCode = excluded.chatCode,
			chatID = excluded.chatID;`)
	if err != nil {
		log.Fatalf("Could not prepare setChatPair query:%v\n", err)
	}
	_, err = query.ExecContext(ctx, cp.UserID, cp.ChatCode, cp.ChatID)
	if err != nil {
		log.Fatalf("Unable to get chat pair:%v\n", err)
	}
}

func registerSecret(ctx context.Context, userID int, secret string) {
	query, err := DB.PrepareContext(ctx, `
	INSERT INTO userbots (userID, secret)
	VALUES ($1, $2)
	ON CONFLICT (userID) DO UPDATE
		SET secret = excluded.secret;`)
	if err != nil {
		log.Fatalf("Could not prepare registerSecret query:%v\n", err)
	}
	_, err = query.ExecContext(ctx, userID, secret)
	if err != nil {
		log.Fatalf("Unable to register secret:%v\n", err)
	}
}

func getSecret(ctx context.Context, userID int) string {
	query, err := DB.PrepareContext(ctx, "SELECT secret FROM userbots WHERE userID = $1;")
	if err != nil {
		log.Fatalf("Could not prepare getSecret query:%v\n", err)
	}
	res := query.QueryRowContext(ctx, userID)

	var secret string
	err = res.Scan(&secret)
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		} else {
			log.Fatalf("Secret scan error:%v\n", err)
		}
	}
	return secret
}
