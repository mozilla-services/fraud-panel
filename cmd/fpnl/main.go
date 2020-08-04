package main

import (
	"fmt"
	"log"
	"os"
)

type config struct {
	OAuthClientID     string
	FraudPanelBaseURL string
}

var cfg = config{
	OAuthClientID:     "771877815709-s9gihlml6lsltg6prj18ddb3jsu74se1.apps.googleusercontent.com",
	FraudPanelBaseURL: "http://127.0.0.1:8000",
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "login":
		sessionID, err := logInUser(cfg.OAuthClientID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(sessionID)
	}
}

func usage() {
	fmt.Println(`
Usage: fpnl COMMAND

A command-line client for fraud-panel

Commands:
  login		Log in to a fraud-panel server using Google Auth

`)
	os.Exit(1)
}
