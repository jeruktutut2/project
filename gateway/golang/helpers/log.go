package helper

import (
	"runtime/debug"
	"strings"
	"time"
)

func PrintLogToTerminal(err error, requestId string) {
	stacktrace := string(debug.Stack())
	stacktrace = strings.ReplaceAll(stacktrace, "\n", "")
	log := `{"logTime": "` + time.Now().String() + `", "app": "project-gateway", "requestId": "` + requestId + `", "stacktrace": "` + stacktrace + `", "error": "` + err.Error() + `"}`
	println(log)
}
