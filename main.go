package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"go.mozilla.org/fraud-panel/db"
	"go.mozilla.org/fraud-panel/mozlog"
)

func main() {
	dbcfg := db.Config{
		Name:     "fpnl",
		User:     "fpnlapi",
		Password: "fpnlapi",
		Host:     "127.0.0.1:5432",
	}
	err := db.Connect(dbcfg)
	if err != nil {
		log.Println(err)
		os.Exit(10)
	}
	go db.Monitor(60 * time.Second)

	listen := "0.0.0.0:8000"
	mux := http.NewServeMux()
	mux.Handle("/auth/pkce", http.HandlerFunc(oAuthPkceExchanger))
	mux.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAccount(w, r)
			return
		case http.MethodPost:
			createAccount(w, r)
			return
		case http.MethodPut:
			updateAccount(w, r)
			return
		case http.MethodDelete:
			deleteAccount(w, r)
			return
		default:
			http.NotFound(w, r)
			return
		}
	})
	mux.Handle("/__version__", http.HandlerFunc(dfVersion))
	mux.Handle("/__heartbeat__", http.HandlerFunc(dfHeartbeat))
	mux.Handle("/__lbheartbeat__", http.HandlerFunc(dfLbHeartbeat))
	server := &http.Server{
		Addr: listen,
		Handler: handleMiddlewares(
			mux,
			setRequestID(),
			setRequestStartTime(),
			setResponseHeaders(),
			logRequest(),
		),
	}
	mozlog.Info("fraud-panel api server listening on %s", listen)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
