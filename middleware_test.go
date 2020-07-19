package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"go.mozilla.org/fraud-panel/mozlog"
)

func TestMiddlewareLogRequest(t *testing.T) {
	oldOut := os.Stdout // keep backup of the real stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	ts := httptest.NewServer(handleMiddlewares(http.HandlerFunc(dfLbHeartbeat), logRequest()))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected http response code 200 but got %d", res.StatusCode)
	}

	w.Close()
	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = oldOut // restoring the real stdout

	var testLog mozlog.Entry
	err = json.Unmarshal(out, &testLog)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", testLog)

	if _, ok := testLog.Fields["msg"]; !ok {
		t.Errorf("expected 'msg' field not found")
	} else {
		if testLog.Fields["msg"] != "request" {
			t.Errorf("'msg' field not set to 'request'")
		}
	}

	if _, ok := testLog.Fields["method"]; !ok {
		t.Errorf("expected 'method' field not found")
	} else {
		if testLog.Fields["method"] != "GET" {
			t.Errorf("'msg' field not set to 'GET'")
		}
	}

	if _, ok := testLog.Fields["url"]; !ok {
		t.Errorf("expected 'url' field not found")
	} else {
		if testLog.Fields["url"] != "/" {
			t.Errorf("'url' field not set to '/'")
		}
	}
}

func TestMiddlewareSetResponseHeaders(t *testing.T) {
	ts := httptest.NewServer(handleMiddlewares(http.HandlerFunc(dfLbHeartbeat), setResponseHeaders()))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected http response code 200 but got %d", res.StatusCode)
	}
	t.Logf("%+v\n", res.Header)
	if _, ok := res.Header["Content-Security-Policy"]; !ok {
		t.Errorf("expected 'Content-Security-Policy' response header not found")
	}
	if _, ok := res.Header["Strict-Transport-Security"]; !ok {
		t.Errorf("expected 'Strict-Transport-Security' response header not found")
	}
}
