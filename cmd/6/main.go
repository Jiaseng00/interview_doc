package main

import (
	"errors"
	"fmt"
	"github.com/natefinch/lumberjack"
	pkgErr "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	ErrorDB    = errors.New("DB connection error")
	ErrorTwo   = errors.New("error two")
	ErrorThree = errors.New("error three")
)

func main() {
	log := logrus.New()
	// Set up file-based log rotation with JSON format
	fileLogger := &lumberjack.Logger{
		Filename:   "service.log", // Log file name
		MaxSize:    100,           // Max size in MB before rotating
		MaxBackups: 30,            // Max number of old log files to keep
		MaxAge:     28,            // Max number of days to retain old log files
		Compress:   true,          // Compress old log files
	}

	// Create a separate logger for file logging with JSON format
	fileLog := logrus.New()
	fileLog.SetFormatter(&logrus.JSONFormatter{})
	fileLog.SetOutput(fileLogger)

	// Create a separate logger for console logging with color
	consoleLog := logrus.New()
	consoleLog.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true, // Forces colorization on the output
		FullTimestamp:   true, // Adds full timestamps
		TimestampFormat: fmt.Sprintf("%s", time.Now().UTC().Format("2006-01-02 15:04:05")),
	})
	consoleLog.SetOutput(os.Stdout)

	// Set log level (Info and above)
	log.SetLevel(logrus.InfoLevel)

	// Log to console and file
	consoleLog.Info("Service started")
	fileLog.Info("Service started")

	consoleLog.Warn("TEST")
	fileLog.Warn("TEST")

	// Log with contextual information
	consoleLog.WithFields(logrus.Fields{
		"service_name": "UserService",
		"request_id":   "12345",
	}).Info("Handling request")

	fileLog.WithFields(logrus.Fields{
		"service_name": "UserService",
		"request_id":   "12345",
	}).Info("Handling request")

	// Simulate an error
	err := someFunction(ErrorTwo)
	if err != nil {
		// Use pkgErr.WithStack to add stack trace information
		wrappedErr := pkgErr.Wrap(ErrorDB, "something wrong in DB")
		// Log the error with stack information
		consoleLog.WithFields(logrus.Fields{
			"error": wrappedErr,
		}).Error("Failed to process request")

		fileLog.WithFields(logrus.Fields{
			"error": wrappedErr,
		}).Error("Failed to process request")
	}

}

func someFunction(err error) error {
	return fmt.Errorf("error : %w", err)
}
