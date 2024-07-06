package helpers

import (
	"runtime/debug"
	"strings"
	"time"
)

func PrintLogToTerminal(err error, requestId string) {
	stacktrace := string(debug.Stack())
	stacktrace = strings.ReplaceAll(stacktrace, "\n", "")
	log := `{"grpcLogTime": "` + time.Now().String() + `", "app": "project-user", "requestId": "` + requestId + `", "stacktrace": "` + stacktrace + `", "error": "` + err.Error() + `"}`
	println(log)
}
