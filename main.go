package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/client"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var list = utils.NewTTLList()

func main() {
	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	router := httprouter.New()

	router.GET("/metrics", Metrics)
	router.GET("/DeployNs", NamespaceDeployment)
	go func() {
		http.ListenAndServe(":80", router)
	}()

	initLogger()
	client.Start(ctx, list)
}
func Metrics(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	promhttp.Handler().ServeHTTP(w, r)
}

func NamespaceDeployment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	queryValues := r.URL.Query()
	namespace := queryValues.Get("namespace")
	timeInMinutesStr := queryValues.Get("time")

	//set default time to 5 minutes if not provided
	if timeInMinutesStr == "" {
		timeInMinutesStr = "2"
	}

	// Parse JSON request body
	// Parse timeInMinutes from string to int
	timeInMinutes, err := strconv.Atoi(timeInMinutesStr)
	if err != nil {
		http.Error(w, "Invalid time value", http.StatusBadRequest)
		return
	}
	list.Add(namespace, time.Duration(timeInMinutes)*time.Minute)

	// Send response
	w.WriteHeader(http.StatusCreated)
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
