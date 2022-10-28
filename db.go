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

	query, err = DB.PrepareContext(ctx, "CREATE TABLE IF NOT EXISTS users (userID integer primary key, secret text);")
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
	query, err := DB.PrepareContext(ctx, `INSERT INTO chatpairs SET userID = $1, chatCode = $2, chatID = $3
										ON CONFLICT (userID) DO UPDATE SET chatID = $3`)
	if err != nil {
		log.Fatalf("Could not prepare setChatPair query:%v\n", err)
	}
	_, err = query.ExecContext(ctx, cp.UserID, cp.ChatCode, cp.ChatID)
	if err != nil {
		log.Fatalf("Unable to get chat pair:%v\n", err)
	}
}

func registerSecret(ctx context.Context, userID string) {
	query, err := DB.PrepareContext(ctx, `INSERT INTO users SET userID = $1, secret = $2
										ON CONFLICT (userID) DO UPDATE SET secret = $2`)
	if err != nil {
		log.Fatalf("Could not prepare setChatPair query:%v\n", err)
	}
	_, err = query.ExecContext(ctx, userID)
	if err != nil {
		log.Fatalf("Unable to get chat pair:%v\n", err)
	}
}
