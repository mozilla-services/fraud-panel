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
		mozlog.Info(err.Error(), nil)
		os.Exit(10)
	}
	go db.Monitor(60 * time.Second)

	listen := "0.0.0.0:8000"
	mux := http.NewServeMux()
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
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
