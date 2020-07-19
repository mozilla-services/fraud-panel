package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDockerflowLbHeartbeat(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(dfLbHeartbeat))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if !bytes.Equal(body, []byte("this is fine")) {
		t.Errorf("invalid response body '%s' from endpoint", body)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected http response code 200 but got %d", res.StatusCode)
	}
}

func TestDockerflowVersion(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(dfVersion))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected http response code 200 but got %d", res.StatusCode)
	}
	versionRaw, err := ioutil.ReadFile("version.json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(versionRaw, body) {
		t.Errorf("response body doesn't match version.json data.\n-- got\n%s\n-- expected\n%s\n", body, versionRaw)
	}
}
