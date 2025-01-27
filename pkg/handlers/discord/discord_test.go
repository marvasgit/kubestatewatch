package discord

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/marvasgit/kubestatewatch/config"
	"github.com/marvasgit/kubestatewatch/pkg/event"
)

func TestInit(t *testing.T) {
	s := &Discord{}
	expectedError := fmt.Errorf(dcErrMsg, "Missing Discord webhook URL")

	var Tests = []struct {
		ms  config.Discord
		err error
	}{
		{config.Discord{WebhookURL: "somepath"}, nil},
		{config.Discord{}, expectedError},
	}

	for _, tt := range Tests {
		c := &config.Config{}
		c.Handler.Discord = tt.ms
		if err := s.Init(c); !reflect.DeepEqual(err, tt.err) {
			t.Fatalf("Init(): %v", err)
		}
	}
}

func TestObjectCreated(t *testing.T) {
	expectedDiscordMsg := DiscordMsg{
		Embeds: []DiscordEmbed{
			{
				Color: dcColors["Normal"],
				Title: "A `pod` in namespace `new` has been `Created`:\n`foo`",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "POST" {
			t.Errorf("expected a POST request for ObjectCreated()")
		}
		decoder := json.NewDecoder(r.Body)
		var c DiscordMsg
		if err := decoder.Decode(&c); err != nil {
			t.Errorf("%v", err)
		}
		if !reflect.DeepEqual(c, expectedDiscordMsg) {
			t.Errorf("expected %v, got %v", expectedDiscordMsg, c)
		}
	}))

	ms := &Discord{DcWebhookURL: ts.URL}
	p := event.StatemonitorEvent{
		Name:      "foo",
		Kind:      "pod",
		Namespace: "new",
		Reason:    "Created",
		Status:    "Normal",
	}

	ms.Handle(p)
}

// Tests ObjectDeleted() by passing v1.Pod
func TestObjectDeleted(t *testing.T) {
	expectedDiscordMsg := DiscordMsg{
		Embeds: []DiscordEmbed{
			{
				Color: dcColors["Danger"],
				Title: "A `pod` in namespace `new` has been `Deleted`:\n`foo`",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "POST" {
			t.Errorf("expected a POST request for ObjectDeleted()")
		}
		decoder := json.NewDecoder(r.Body)
		var c DiscordMsg
		if err := decoder.Decode(&c); err != nil {
			t.Errorf("%v", err)
		}
		if !reflect.DeepEqual(c, expectedDiscordMsg) {
			t.Errorf("expected %v, got %v", expectedDiscordMsg, c)
		}
	}))

	ms := &Discord{DcWebhookURL: ts.URL}

	p := event.StatemonitorEvent{
		Name:      "foo",
		Namespace: "new",
		Kind:      "pod",
		Reason:    "Deleted",
		Status:    "Danger",
	}

	ms.Handle(p)
}

// Tests ObjectUpdated() by passing v1.Pod
func TestObjectUpdated(t *testing.T) {
	expectedDiscordMsg := DiscordMsg{
		Embeds: []DiscordEmbed{
			{
				Color: dcColors["Warning"],
				Title: "A `pod` in namespace `new` has been `Updated`:\n`foo`",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "POST" {
			t.Errorf("expected a POST request for ObjectUpdated()")
		}
		decoder := json.NewDecoder(r.Body)
		var c DiscordMsg
		if err := decoder.Decode(&c); err != nil {
			t.Errorf("%v", err)
		}
		if !reflect.DeepEqual(c, expectedDiscordMsg) {
			t.Errorf("expected %v, got %v", expectedDiscordMsg, c)
		}
	}))

	ms := &Discord{DcWebhookURL: ts.URL}

	oldP := event.StatemonitorEvent{
		Name:      "foo",
		Namespace: "new",
		Kind:      "pod",
		Reason:    "Updated",
		Status:    "Warning",
	}

	newP := event.StatemonitorEvent{
		Name:      "foo-new",
		Namespace: "new",
		Kind:      "pod",
		Reason:    "Updated",
		Status:    "Warning",
	}
	_ = newP

	ms.Handle(oldP)
}
