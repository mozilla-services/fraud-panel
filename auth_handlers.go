package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"go.mozilla.org/fraud-panel/db"
	"go.mozilla.org/fraud-panel/errors"
	"go.mozilla.org/fraud-panel/mozlog"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

var oauthCfg = &Config{
	ClientID:     "771877815709-s9gihlml6lsltg6prj18ddb3jsu74se1.apps.googleusercontent.com",
	ClientSecret: "p5pfrMJgtCTJWC4LjzrZoAkE",
	RedirectURL:  "http://127.0.0.1:21337/",
}

type oAuthGoogleError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type Token struct {
	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken string `json:"access_token"`

	// TokenType is the type of token.
	// The Type method returns either this or "Bearer", the default.
	TokenType string `json:"token_type,omitempty"`

	// RefreshToken is a token that's used by the application
	// (as opposed to the user) to refresh the access token
	// if it expires.
	RefreshToken string `json:"refresh_token,omitempty"`

	// Expiry is the optional expiration time of the access token.
	//
	// If zero, TokenSource implementations will reuse the same
	// token forever and RefreshToken or equivalent
	// mechanisms for that TokenSource will not be used.
	Expiry time.Time `json:"expiry,omitempty"`
}

type oAuthGoogleIdentity struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	Hd            string `json:"hd"`
}

// oAuthPkceExchanger receives a code verifier and authorization token from a client,
// queries google for an access token
func oAuthPkceExchanger(w http.ResponseWriter, r *http.Request) {
	rid := getRequestID(r)
	if r.Method != http.MethodPost {
		errors.ErrMustBePostRequest.HTTPError(w, r, rid)
		return
	}
	codeVerifier := r.FormValue("code_verifier")
	if codeVerifier == "" {
		errors.ErrPKCECodeVerifierMissing.HTTPError(w, r, rid)
		return
	}
	authorizationCode := r.FormValue("authorization_code")
	if codeVerifier == "" {
		errors.ErrPKCEAuthorizationCodeMissing.HTTPError(w, r, rid)
		return
	}
	resp, err := http.PostForm("https://oauth2.googleapis.com/token",
		url.Values{
			"client_id":     {oauthCfg.ClientID},
			"client_secret": {oauthCfg.ClientSecret},
			"code_verifier": {codeVerifier},
			"code":          {authorizationCode},
			"redirect_uri":  {oauthCfg.RedirectURL},
			"grant_type":    {"authorization_code"},
		})
	if err != nil {
		errors.ErrPKCEAuthorizationFailed.HTTPError(w, r, rid)
		mozlog.Event(
			fmt.Sprintf("googleapi responded with http status %d %s",
				resp.StatusCode, resp.Status),
			mozlog.Fields{"rid": rid},
		)
		return
	}
	// process the response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errors.ErrPKCEAuthorizationFailed.HTTPError(w, r, rid)
		mozlog.Event(fmt.Sprintf("googleapi responded with non-parseable body: %v", err),
			mozlog.Fields{"rid": rid},
		)
		return
	}
	if resp.StatusCode != http.StatusOK {
		errors.ErrPKCEAuthorizationFailed.HTTPError(w, r, rid)
		//unmarshal response body as an error
		var gErr oAuthGoogleError
		json.Unmarshal(body, &gErr)
		mozlog.Event(
			fmt.Sprintf("googleapi responded with http status %d %s",
				resp.StatusCode, resp.Status),
			mozlog.Fields{"rid": rid, "error": gErr.Error, "error_description": gErr.ErrorDescription},
		)
		return
	}
	var token Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		errors.ErrPKCEAuthorizationFailed.HTTPError(w, r, rid)
		mozlog.Event(fmt.Sprintf("googleapi responded with non-parseable body: %v", err),
			mozlog.Fields{"rid": rid},
		)
		return
	}
	sessionID, err := makeUserSession(token.AccessToken)
	if err != nil {
		errors.ErrMakingUserSession.HTTPError(w, r, rid)
		mozlog.Event(fmt.Sprintf("failed to make a session for the user: %v", err),
			mozlog.Fields{"rid": rid},
		)
		return
	}
	c := http.Cookie{
		Name:     "fpnl-session-id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   86400,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &c)
	w.Write([]byte("Session created. You are logged in."))
	mozlog.Event("user session created successfully",
		mozlog.Fields{"rid": rid})
}

// makeUserSession retrieves the identify of a user from Google using an access token
// then creates a local session ID in database and returns it to the user
func makeUserSession(accessToken string) (sessionID string, err error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo?alt=json", nil)
	if err != nil {
		mozlog.Info("failed to create new request to retrieve user information: %v", err)
		return "", fmt.Errorf("%w: failed to create new request", errors.ErrOAuthUserInfoFailed)
	}
	req.Header.Add("Authorization", `Bearer `+accessToken)
	var client = new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		mozlog.Info("failed to retrieve user information from google: %v", err)
		return "", fmt.Errorf("%w: failed to retrieve user information from google", errors.ErrOAuthUserInfoFailed)
	}
	// process the response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		mozlog.Info("failed to read response from google: %v", err)
		return "", fmt.Errorf("%w: failed to read response from google", errors.ErrOAuthUserInfoFailed)
	}
	if resp.StatusCode != http.StatusOK {
		//unmarshal response body as an error
		var gErr oAuthGoogleError
		json.Unmarshal(body, &gErr)
		mozlog.Event(
			fmt.Sprintf("google responded with http status %d %s", resp.StatusCode, resp.Status),
			mozlog.Fields{"error": gErr.Error, "error_description": gErr.ErrorDescription},
		)
		return "", fmt.Errorf("%w: google responded with http status %d %s",
			resp.StatusCode, resp.Status, errors.ErrOAuthUserInfoFailed)
	}
	var oagi oAuthGoogleIdentity
	err = json.Unmarshal(body, &oagi)
	if err != nil {
		mozlog.Info("googleapi responded with non-parseable body: %v", err)
		return "", fmt.Errorf("%w: googleapi responded with non-parseable body",
			errors.ErrOAuthUserInfoFailed)
	}
	mozlog.Event("retrieved user details from google",
		mozlog.Fields{"name": oagi.Name, "email": oagi.Email})
	return db.UserLogIn(oagi.Name, oagi.Email)
}
