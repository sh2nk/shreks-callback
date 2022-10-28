package main

import (
	"context"
	"log"
	"net/http"

	iris "github.com/sh2nk/shreks-callback/iris-callback-api"
)

func Auth(ctx context.Context, next http.Handler) http.Handler {

}

func Callback(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signal, err := iris.UnmarshalSignal(r.Body)
		if err != nil {
			log.Printf("Unmarshal signal error: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}

		log.Printf("Got Iris %s signal!\n%v\n", signal.Method, signal)

		switch signal.Method {
		case "addUser":
			OnAddUser(ctx, w, signal)
		case "subscribeSignals":
			OnSubscribeSignals(ctx, w, signal)
		}

		next.ServeHTTP(w, r)
	})
}

func OK(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
