package mozlog //import "go.mozilla.org/fraud-panel/mozlog"

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

// MozLogger implements the io.Writer interface
type MozLogger struct {
	Output io.Writer
	Name   string
}

var logger = &MozLogger{
	Output: os.Stdout,
	Name:   "Application",
}

var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		log.Printf("Can't resolve hostname: %v", err)
	}

	log.SetOutput(logger)
	log.SetFlags(log.Lmsgprefix)
}

// Write converts the log to AppLog
func (m *MozLogger) Write(l []byte) (int, error) {
	var f Fields
	err := json.Unmarshal(l, &f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 0, err
	}
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
		return 0, err
	}

	_, err = m.Output.Write(append(out, '\n'))
	return len(l), err
}

// Info logs a set of fields
func Info(msg string, f Fields) error {
	if f == nil {
		f = make(Fields)
	}
	f["msg"] = msg
	jsonFields, err := json.Marshal(f)
	if err != nil {
		return err
	}
	log.Printf(string(jsonFields))
	return nil
}
