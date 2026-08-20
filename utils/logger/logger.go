package logger

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Log is a global package-level instance of Logrus Logger.
var Log *logrus.Logger

// InitLogger initializes the global Log logger with colorized formatted text output.
func InitLogger() {
	Log = logrus.New()

	// Use TextFormatter with forced colors for formatted console output
	Log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableQuote:    true,
		PadLevelText:    true,
	})

	// Output to stdout instead of standard stderr
	Log.SetOutput(os.Stdout)

	// Set logging level
	Log.SetLevel(logrus.InfoLevel)
}

// RequestLoggerMiddleware returns a Gin middleware that logs the start and end of HTTP requests with clear delimiters.
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		if Log != nil {
			Log.Infof("--- [START] API Request: %s %s ---", method, path)
		}

		c.Next()

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()

		if Log != nil {
			Log.Infof("---- [END] API Request: %s %s | Status: %d | Latency: %v ----", method, path, statusCode, latency)
		}
	}
}
