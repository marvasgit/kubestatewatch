package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/marvasgit/KubeStateWatch/pkg/client"
	"github.com/marvasgit/KubeStateWatch/pkg/utils"
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
	router.PUT("/deploy/:namespace/:duration", namespaceDeployment)
	router.PUT("/deploy/:namespace", namespaceDeployment)
	router.DELETE("/deploy/:namespace", deletenamespaceDeployment)
	router.POST("/reset", reset)
	go func() {
		http.ListenAndServe(":80", router)
	}()

	initLogger()
	client.Start(ctx, list)
}
func Metrics(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	promhttp.Handler().ServeHTTP(w, r)
}
func deletenamespaceDeployment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	namespace := ps.ByName("namespace")

	list.Remove(namespace)
	response := fmt.Sprintf("Namespace -%s was removed ", namespace)
	w.Write([]byte(response))
}
func reset(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list.Reset()
	w.Write([]byte("Deployment List was reset "))
}

func namespaceDeployment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	namespace := ps.ByName("namespace")
	durationString := ps.ByName("duration")

	//set default time to 5 minutes if not provided
	if durationString == "" {
		durationString = "2"
	}

	// Parse JSON request body
	// Parse durationInMinutes from string to int
	durationInMinutes, err := strconv.Atoi(durationString)
	if err != nil {
		http.Error(w, "Invalid time value", http.StatusBadRequest)
		return
	}
	er := list.Add(namespace, time.Duration(durationInMinutes)*time.Minute)
	if er != nil {
		http.Error(w, er.Error(), http.StatusBadRequest)
		return
	}
	response := fmt.Sprintf("Namespace -%s added to deployment list", namespace)
	// change status code to 201
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
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
