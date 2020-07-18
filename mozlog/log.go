package mozlog //import "go.mozilla.org/fraud-panel/mozlog"

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry implements a log entry in Mozilla logging standard
type Entry struct {
	Timestamp  int64
	Time       string
	Type       string
	Logger     string
	Hostname   string `json:",omitempty"`
	EnvVersion string
	Pid        int `json:",omitempty"`
	Severity   int `json:",omitempty"`
	Fields     Fields
}

// Fields adds a map of fields to a log entry
type Fields map[string]interface{}

var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't resolve hostname: %v", err)
		os.Exit(97)
	}
}

// Event writes an Entry to stdout
func Event(msg string, f Fields) {
	if f == nil {
		f = make(Fields)
	}
	f["msg"] = msg
	log := Entry{
		Timestamp:  time.Now().UnixNano(),
		Time:       time.Now().Format(time.RFC3339),
		Type:       "app.log",
		Logger:     "fraud-panel",
		Hostname:   hostname,
		EnvVersion: "2.0",
		Pid:        os.Getpid(),
		Fields:     f,
	}

	out, err := json.Marshal(log)
	if err != nil {
		// Need someway to notify that this happened.
		fmt.Fprintln(os.Stderr, err)
		os.Exit(13)
	}
	fmt.Fprintf(os.Stdout, "%s\n", out)
}

// Info logs a simple message
func Info(format string, a ...interface{}) {
	Event(fmt.Sprintf(format, a...), nil)
}
