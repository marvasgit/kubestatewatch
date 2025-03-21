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

package flock

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/marvasgit/kubestatewatch/config"
	"github.com/marvasgit/kubestatewatch/pkg/event"
)

var flockColors = map[string]string{
	"Normal":  "#00FF00",
	"Warning": "#FFFF00",
	"Danger":  "#FF0000",
}

var flockErrMsg = `
%s

You need to set Flock url for Flock notify,
using "--url/-u" or using environment variables:

export KW_FLOCK_URL=flock_url

Command line flags will override environment variables

`

// Flock handler implements handler.Handler interface,
// Notify event to Flock channel
type Flock struct {
	Url string
}

// FlockMessage struct
type FlockMessage struct {
	Notification string                    `json:"notification"`
	Text         string                    `json:"text"`
	Attachements []FlockMessageAttachement `json:"attachments"`
}

// FlockMessageAttachement struct
type FlockMessageAttachement struct {
	Title string                       `json:"title"`
	Color string                       `json:"color"`
	Views FlockMessageAttachementViews `json:"views"`
}

// FlockMessageAttachementViews struct
type FlockMessageAttachementViews struct {
	Flockml string `json:"flockml"`
}

// Init prepares Flock configuration
func (f *Flock) Init(c *config.Config) error {
	url := c.Handler.Flock.Url

	if url == "" {
		url = os.Getenv("KW_FLOCK_URL")
	}

	f.Url = url

	return checkMissingFlockVars(f)
}

// Handle handles an event.
func (f *Flock) Handle(e event.StatemonitorEvent) {
	flockMessage := prepareFlockMessage(e, f)

	err := postMessage(f.Url, flockMessage)
	if err != nil {
		logrus.Errorf("Failed to send flock message: %v\n", err)
		return
	}

	logrus.Printf("Message successfully sent to channel %s at %s", f.Url, time.Now())
}

func checkMissingFlockVars(s *Flock) error {
	if s.Url == "" {
		return fmt.Errorf(flockErrMsg, "Missing Flock url")
	}

	return nil
}

func prepareFlockMessage(e event.StatemonitorEvent, f *Flock) *FlockMessage {
	return &FlockMessage{
		Text:         "statemonitor Alert",
		Notification: "statemonitor Alert",
		Attachements: []FlockMessageAttachement{
			{
				Title: e.Message(),
				Color: flockColors[e.Status],
			},
		},
	}
}

func postMessage(url string, flockMessage *FlockMessage) error {
	message, err := json.Marshal(flockMessage)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
