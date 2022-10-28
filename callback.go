package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	iris "github.com/sh2nk/shreks-callback/iris-callback-api"
)

type contextKey string

const signalKey contextKey = "signal"

func HandleSignal(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if incoming request is valid Iris signal and try to parse it
		signal, err := iris.UnmarshalSignal(r.Body)
		if err != nil {
			log.Printf("Unmarshal signal error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Some logging action
		log.Printf("Got Iris %s signal!\n%v\n", signal.Method, signal)

		// Store parsed signal data in context
		ctx := context.WithValue(r.Context(), signalKey, signal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get context
		ctx := r.Context()

		// Get signal value from conetxt
		signal, ok := ctx.Value(signalKey).(iris.IrisSignal)
		if !ok {
			err := errors.New("signal context error")
			log.Printf("Can't perform auth: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check user secret
		if signal.Secret != getSecret(ctx, signal.UserID) {
			err := errors.New("Auth error: invalid secret")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Callback(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get context
		ctx := r.Context()

		// Get signal value from conetxt
		signal, ok := ctx.Value(signalKey).(iris.IrisSignal)
		if !ok {
			err := errors.New("signal context error")
			log.Printf("Callback falied: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		switch signal.Method {
		case "banExpired":
			OnBanExpired(ctx, w)
		case "addUser":
			OnAddUser(ctx, w)
		case "subscribeSignals":
			OnSubscribeSignals(ctx, w)
		case "deleteMessages":
			OnDeleteMessages(ctx, w)
		case "deleteMessagesFromUser":
			OnDeleteMessagesFromUser(ctx, w)
		}

		next.ServeHTTP(w, r)
	})
}

func OK() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
}
