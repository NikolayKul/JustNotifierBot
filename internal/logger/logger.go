package logger

import (
	stdlog "log"
	"os"
	"strings"
)

var logger = createLogger()

// Println calls to Logger
func Println(v ...interface{}) {
	logger.Println(v...)
}

// Printf calls to Logger
func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

// Fatal calls to Logger
func Fatal(err error) {
	logger.Fatal(err)
}

type emptyOutput struct{}

func (o *emptyOutput) Write(p []byte) (n int, err error) {
	return 0, nil
}

func createLogger() *stdlog.Logger {
	defaultLogger := stdlog.New(os.Stdout, "", 0)
	if !isInDebugMode() {
		defaultLogger.SetOutput(&emptyOutput{})
	}
	return defaultLogger
}

func isInDebugMode() bool {
	value, exists := os.LookupEnv("JUST_NOTIFIER_BOT_IS_DEBUG_MODE")
	return exists && (value == "1" || strings.ToLower(value) == "true")
}
