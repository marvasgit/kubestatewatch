package main

import (
	"context"
	"encoding/json"
	"fmt"
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
var nschan = make(chan httpInput)

func main() {
	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	router := httprouter.New()

	router.GET("/metrics", Metrics)
	router.GET("/DeployNs", NamespaceDeployment)

	//log.Fatal(http.ListenAndServe(":8080", router))

	// go func() {
	// 	http.Handle("/metrics", promhttp.Handler())
	// 	http.ListenAndServe(":2112", nil)
	// }()
	// go func() {
	// 	http.Handle("/add", http.HandlerFunc(addItemHandler))
	// 	http.ListenAndServe(":2113", nil)
	// }()
	initLogger()
	client.Start(ctx, nschan)
}
func Metrics(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	promhttp.Handler().ServeHTTP(w, r)
}

func NamespaceDeployment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	namespace := ps.ByName("namespace")
	timeInMinutesStr := ps.ByName("time")

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

func addItemHandler(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	decoder := json.NewDecoder(r.Body)
	item := httpInput
	err := decoder.Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add item to the list
	list.Add(item.Value, time.Duration(item.ExpiresAfter)*time.Minute)

	// Send response
	w.WriteHeader(http.StatusCreated)
}

type httpInput struct {
	Namespace    string `json:"namespace"`
	ExpiresAfter int8   `json:"expires_after_minutes"`
}

func StartUpdateRoutine(ctx context.Context, resource *TTLList, updateCh <-chan httpInput) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case value := <-updateCh:

				if err := resource.Add(value.Namespace, time.Duration(value.ExpiresAfter)*time.Minute); err != nil {
					fmt.Printf("Error setting value: %v\n", err)
				}
			}
		}
	}()
}
