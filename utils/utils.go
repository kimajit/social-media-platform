package utils

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// InitializeLogger sets up the logger.
func InitializeLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	return logger
}

// RespondWithError sends an error response.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	http.Error(w, message, code)
}
