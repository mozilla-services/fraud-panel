package errors

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"go.mozilla.org/fraud-panel/mozlog"
)

// Code defines an error code
type Code uint16

const (
	// ErrMustBePostRequest request must use HTTP Post method
	ErrMustBePostRequest Code = iota + 1000
)

const (
	// ErrAccountIDMissing missing account id
	ErrAccountIDMissing Code = iota + 2000
	// ErrPKCECodeVerifierMissing missing pkce code verifier
	ErrPKCECodeVerifierMissing
	// ErrPKCEAuthorizationCodeMissing missing pce authorization code
	ErrPKCEAuthorizationCodeMissing
	// ErrPKCEAuthorizationFailed failed to exchange authorization code for access token
	ErrPKCEAuthorizationFailed
	// ErrOAuthUserInfoFailed failed to retrieve user info
	ErrOAuthUserInfoFailed
	// ErrMakingUserSession failed to create a session for the user
	ErrMakingUserSession
)

func (c Code) String() string {
	switch c {
	case ErrMustBePostRequest:
		return "http POST method required"
	case ErrAccountIDMissing:
		return "missing account id"
	case ErrPKCECodeVerifierMissing:
		return "pkce code verifier not provided in form values"
	case ErrPKCEAuthorizationCodeMissing:
		return "pkce authorization code not provided in form values"
	case ErrPKCEAuthorizationFailed:
		return "pkce authorization failed"
	case ErrOAuthUserInfoFailed:
		return "failed to retrieve user information from google"
	case ErrMakingUserSession:
		return "failed to create a user session"
	}
	return ""
}

// StatusCode returns an appropriate HTTP status code for a given error
func (c Code) StatusCode() int {
	switch c {
	case ErrMustBePostRequest:
		return http.StatusMethodNotAllowed
	case ErrAccountIDMissing, ErrPKCECodeVerifierMissing, ErrPKCEAuthorizationCodeMissing:
		return http.StatusBadRequest
	case ErrPKCEAuthorizationFailed:
		return http.StatusUnauthorized
	case ErrOAuthUserInfoFailed, ErrMakingUserSession:
		return http.StatusInternalServerError
	}
	return 0
}

// HTTPError replies to the request with the specified error code and message
func (c Code) HTTPError(w http.ResponseWriter, r *http.Request, rid string) {
	mozlog.Event(c.String(), mozlog.Fields{"http_status": c.StatusCode(), "error_code": c, "rid": rid})
	// when nginx is in front of go, nginx requires that the entire
	// request body is read before writing a response.
	// https://github.com/golang/go/issues/15789
	if r.Body != nil {
		io.Copy(ioutil.Discard, r.Body)
		r.Body.Close()
	}
	http.Error(w, fmt.Sprintf("%s\ncode: %d\nrequest id: %s\n", c, c, rid), c.StatusCode())
	return
}
