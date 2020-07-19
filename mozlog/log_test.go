package mozlog

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestEvent(t *testing.T) {
	oldOut := os.Stdout // keep backup of the real stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	Event("this is a test", Fields{"foo": "bar"})

	w.Close()
	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = oldOut // restoring the real stdout

	var testLog Entry
	err = json.Unmarshal(out, &testLog)
	if err != nil {
		t.Fatal(err)
	}
	if testLog.Fields["msg"] != "this is a test" {
		t.Errorf("expected to find log msg 'this is a test' but found '%s'", testLog.Fields["msg"])
	}
	if _, ok := testLog.Fields["foo"]; !ok {
		t.Errorf("expected to find field key 'foo' but didn't")
	} else {
		if testLog.Fields["foo"] != "bar" {
			t.Errorf("expected to find field key 'foo' with value 'bar' but found '%s'", testLog.Fields["foo"])
		}
	}
	if testLog.Timestamp == 0 {
		t.Errorf("invalid empty timestamp")
	}
	if testLog.Time == "" {
		t.Errorf("invalid empty time")
	}
	if testLog.Hostname == "" {
		t.Errorf("invalid empty hostname")
	}
	testHostName, err := os.Hostname()
	if err != nil {
		t.Fatal(err)
	}
	if testLog.Hostname != testHostName {
		t.Errorf("hostname mismatch. expected %q got %q", testHostName, testLog.Hostname)
	}
}
