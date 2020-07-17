package main

import (
	"log"
	"net/http"
)

func main() {
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
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
