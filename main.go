/*
Copyright 2016 Skippbox, Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"net/http"
	"os"

	"github.com/marvasgit/kubernetes-diffwatcher/pkg/client"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	client.RunWithConfig()
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
