package main

import (
	"context"
	"log"
)

func createTables(ctx context.Context) {
	query, err := DB.PrepareContext(ctx, "CREATE TABLE IF NOT EXISTS userbots (userID integer primary key, chat test);")
	if err != nil {
		log.Fatalf("Could not prepare db create query:%v\n", err)
	}
	_, err = query.ExecContext(ctx)
	if err != nil {
		log.Fatalf("Unable to create database: %v\n", err)
	}
	defer query.Close()
}
