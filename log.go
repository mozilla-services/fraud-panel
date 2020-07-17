package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const appLogName = "fraud-panel"

var hostname = "unknown"

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		panic(fmt.Sprintf("cannot determine hostname: %v", err))
	}
}

// logRequest is a middleware that writes details about each HTTP request processed
// but the various handlers. It is executed last to capture signing logs as well.
func logRequest() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
			// attempt to retrieve a signing registry entry for this request
			// from the global sr.entry map, using mutexes
			rid := getRequestID(r)
			// calculate the processing time
			t1 := getRequestStartTime(r)
			procTs := time.Now().Sub(t1)
			fmt.Printf(`{"Timestamp":%d,"Time":%q,"Type":"app.log","Logger":%q,`+
				`"Hostname":%q,"EnvVersion":"2.0","Pid":%d,"Severity":6,"Fields":{`+
				`"method":%q,"msg":"request","proto":%q,"remoteAddress":%q,`+
				`"remoteAddressChain":"[%s]","rid":%q,"t":%d,"ua":%q,"url":%q}}
`,
				time.Now().UnixNano(), time.Now().Format(time.RFC3339), appLogName,
				hostname, os.Getpid(), r.Method, r.Proto, r.RemoteAddr,
				r.Header.Get("X-Forwarded-For"), rid, procTs/time.Millisecond,
				r.UserAgent(), r.URL.String())
		})
	}
}

// setRequestID is a middleware the generates a random ID for each request processed
// by the HTTP server. The request ID is added to the request context and used to
// track various information and correlate logs.
func setRequestID() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rid := make([]rune, 16)
			letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
			for i := range rid {
				rid[i] = letters[rand.Intn(len(letters))]
			}

			h.ServeHTTP(w, addToContext(r, contextKeyRequestID, string(rid)))
		})
	}
}

// setRequestStartTime is a middleware that stores a timestamp of the time a request entering
// the middleware, to calculate processing time later on
func setRequestStartTime() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, addToContext(r, contextKeyRequestStartTime, time.Now()))
		})
	}
}
