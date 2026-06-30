package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Log is a global package-level instance of Logrus Logger.
var Log *logrus.Logger

// InitLogger initializes the global Log logger.
func InitLogger() {
	Log = logrus.New()

	// Use JSON format for structured production logs
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Output to stdout instead of standard stderr
	Log.SetOutput(os.Stdout)

	// Set logging level
	Log.SetLevel(logrus.InfoLevel)
}
