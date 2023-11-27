package main

import (
	"net/http"
	"os"

	"github.com/marvasgit/kubernetes-statemonitor/pkg/client"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()
	initLogger()
	client.Start()
}

func initLogger() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		logrus.Printf("Custom log level: %s", logLevel)
		parsedLevel, err := logrus.ParseLevel(logLevel)
		if err == nil {
			logrus.Printf("Setting custom log level to: %s", logLevel)
			logrus.SetLevel(parsedLevel)
		} else {
			logrus.Errorf("Illegal custom log level: %s. Ignoring custom log level", logLevel)
		}
	}
	logFormatter := os.Getenv("LOG_FORMATTER")
	if logFormatter != "" {
		logrus.Printf("Custom log formatter: %s", logFormatter)
		if logFormatter == "json" {
			logrus.Printf("Setting custom log formatter to: %s", logFormatter)
			logrus.SetFormatter(new(logrus.JSONFormatter))
		} else {
			logrus.Errorf("Illegal custom log formatter: %s. Ignoring custom log formatter", logFormatter)
		}
	}
}
