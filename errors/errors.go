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
	// ErrAccountIDMissing missing account id
	ErrAccountIDMissing Code = iota + 1000
)

func (c Code) String() string {
	switch c {
	case ErrAccountIDMissing:
		return "missing account id"
	}
	return ""
}

// StatusCode returns an appropriate HTTP status code for a given error
func (c Code) StatusCode() int {
	switch c {
	case ErrAccountIDMissing:
		return http.StatusBadRequest
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
	http.Error(w, fmt.Sprintf("%s\ncode: %d", c, c), c.StatusCode())
	return
}
