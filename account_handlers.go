package main

import (
	"net/http"

	"go.mozilla.org/fraud-panel/errors"
)

func getAccount(w http.ResponseWriter, r *http.Request) {
	rid := getRequestID(r)
	accountID := r.FormValue("id")
	if accountID == "" {
		errors.ErrAccountIDMissing.HTTPError(w, r, rid)
		return
	}
}

func createAccount(w http.ResponseWriter, r *http.Request) {
}

func updateAccount(w http.ResponseWriter, r *http.Request) {
	rid := getRequestID(r)
	accountID := r.FormValue("id")
	if accountID == "" {
		errors.ErrAccountIDMissing.HTTPError(w, r, rid)
		return
	}
}

func deleteAccount(w http.ResponseWriter, r *http.Request) {
	rid := getRequestID(r)
	accountID := r.FormValue("id")
	if accountID == "" {
		errors.ErrAccountIDMissing.HTTPError(w, r, rid)
		return
	}
}
