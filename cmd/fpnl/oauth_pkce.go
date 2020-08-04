package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"go.mozilla.org/fraud-panel/cmd/fpnl/open"
)

type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

// authorizeUser implements the PKCE OAuth2 flow.
func logInUser(clientID string) (sessionID string, err error) {
	codeVerifier := generateCode()
	open.Open(fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth"+
			"?scope=email%%20profile"+
			"&response_type=code"+
			"&client_id=%s"+
			"&code_challenge=%s"+
			"&code_challenge_method=plain"+
			"&redirect_uri=http://127.0.0.1:21337/",
		clientID, codeVerifier))
	return authHandler(codeVerifier)
}

func generateCode() string {
	msg := make([]byte, 64)
	rand.Read(msg)
	return base64.URLEncoding.EncodeToString(msg)
}

func authHandler(codeVerifier string) (sessionID string, err error) {
	server := &http.Server{Addr: ":21337"}
	// define a handler that will get the authorization code,
	// call the token endpoint, and close the HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("received %s redirect from google", r.Method)
		authorizationCode := r.URL.Query().Get("code")
		if authorizationCode == "" {
			errMsg := fmt.Sprintf("failed to obtain an authorization code from google")
			log.Println(errMsg)
			w.Write([]byte(errMsg))
			server.Close()
			return
		}
		// exchange the authorization code and the code verifier for a session
		sessionID, err = exchangeCodeForSession(codeVerifier, authorizationCode)
		if err != nil {
			errMsg := fmt.Sprintf("failed to obtain session from fraud-panel server: %v", err)
			log.Println(errMsg)
			w.Write([]byte(errMsg))
			server.Close()
			return
		}
		// return an indication of success to the caller
		w.Write([]byte(`<html>
			<body>
				<h1>Login successful!</h1>
				<h2>You can close this browsing window and return to the command line.</h2>
			</body>
		</html>`))
		go server.Close()
	})
	err = server.ListenAndServe()
	if err == http.ErrServerClosed {
		// in this situation, closing the server is not an error
		err = nil
	}
	return
}

// exchangeCodeForSession trades the authorization code retrieved from the first OAuth2 leg
// for a session cookie on the fraud panel
func exchangeCodeForSession(codeVerifier string, authorizationCode string) (sessionID string, err error) {
	resp, err := http.PostForm(cfg.FraudPanelBaseURL+"/auth/pkce",
		url.Values{
			"code_verifier":      {codeVerifier},
			"authorization_code": {authorizationCode},
		})
	if err != nil {
		return
	}
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "fpnl-session-id" {
			sessionID = cookie.Value
		}
	}
	return
}
