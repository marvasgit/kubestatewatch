package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/marvasgit/kubernetes-diffwatcher/config"
	"github.com/marvasgit/kubernetes-diffwatcher/pkg/event"
	"github.com/marvasgit/kubernetes-diffwatcher/pkg/handlers"
	"github.com/sirupsen/logrus"
)

var dcErrMsg = `
%s

You need to set the MS teams webhook URL,
using --webhookURL, or using environment variables:

export KW_DISCORD_WEBHOOKURL=webhook_url

Command line flags will override environment variables

`

var dcColors = map[string]int{
	"Normal":  8311585,
	"Warning": 16312092,
	"Danger":  13632027,
}

type Discord struct {
	DcWebhookURL string
}

type DiscordMsg struct {
	Embeds []DiscordEmbed `json:"embeds"`
}

type DiscordEmbed struct {
	Color int    `json:"color"`
	Title string `json:"title"`
}

var _ handlers.Handler = &Discord{}

func (dc *Discord) Init(c *config.Config) error {
	webhookURL := c.Handler.Discord.WebhookURL

	if webhookURL == "" {
		webhookURL = os.Getenv("KW_DISCORD_WEBHOOKURL")
	}

	if webhookURL == "" {
		return fmt.Errorf(dcErrMsg, "Missing Discord webhook URL")
	}

	dc.DcWebhookURL = webhookURL
	return nil
}

func (dc *Discord) Handle(e event.DiffWatchEvent) {
	msg := &DiscordMsg{}

	var embed DiscordEmbed
	embed.Color = dcColors[e.Status]
	embed.Title = e.Message()

	msg.Embeds = append(msg.Embeds, embed)

	_, err := sendMessage(dc, msg)
	if err != nil {
		logrus.Printf("%s\n", err)
		return
	}

	logrus.Printf("Message successfully sent to Discord")
}

func sendMessage(dc *Discord, discordMsg *DiscordMsg) (*http.Response, error) {
	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(discordMsg); err != nil {
		return nil, fmt.Errorf("Failed encoding message: %v", err)
	}

	res, err := http.Post(dc.DcWebhookURL, "application/json", buffer)
	if err != nil {
		return nil, fmt.Errorf("Failed sending to webhook url %s. Got the error: %v", dc.DcWebhookURL, err)
	}

	if res.StatusCode != http.StatusOK {
		resMessage, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Failed reading Discord http response: %v", err)
		}
		return nil, fmt.Errorf("Failed sending to Discord channel. Discord http response: %s, %s",
			res.Status, string(resMessage))
	}

	if err := res.Body.Close(); err != nil {
		return nil, err
	}

	return res, nil
}
