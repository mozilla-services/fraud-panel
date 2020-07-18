package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.mozilla.org/fraud-panel/db"
	"go.mozilla.org/fraud-panel/mozlog"
)

func dfHeartbeat(w http.ResponseWriter, r *http.Request) {
	var (
		dbAccessible       bool
		dbCheckTimeout     time.Duration = 10 * time.Second
		dbHeartbeatStartTs time.Time     = time.Now()
	)
	dbCheckCtx, dbCancel := context.WithTimeout(r.Context(), dbCheckTimeout)
	defer dbCancel()
	err := db.CheckConnectionContext(dbCheckCtx)
	if err == nil {
		mozlog.Info("db heartbeat completed successfully", mozlog.Fields{
			"rid":     getRequestID(r),
			"t":       int32(time.Since(dbHeartbeatStartTs) / time.Millisecond),
			"timeout": dbCheckTimeout.String(),
		})
		dbAccessible = true
	} else {
		mozlog.Info("error checking DB connection: %s"+err.Error(), nil)
		dbAccessible = false
	}
	w.Write([]byte(fmt.Sprintf(`{"dbAccessible": %t}`, dbAccessible)))
}

func dfLbHeartbeat(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is fine"))
}

func dfVersion(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("version.json")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load version.json: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
