package client

import (
	"context"
	"os"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/marvasgit/kubernetes-statemonitor/config"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/controller"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/handlers"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/handlers/cloudevent"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/handlers/flock"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/handlers/hipchat"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/handlers/lark"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/handlers/mattermost"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/handlers/msteam"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/handlers/slack"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/handlers/slackwebhook"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/handlers/smtpClient"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/handlers/webhook"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/utils"
	"github.com/sirupsen/logrus"
)

func Start(ctx context.Context, list *utils.TTLList) {

	conf := loadConfig()
	handlers := parseEventHandler(&conf)
	controller.Start(&conf, handlers, list)
}

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var k = koanf.New(".")

func loadConfig() config.Config {
	// Load JSON config.
	//read envVariable IsLOCAL
	isLocal := os.Getenv("IsLOCAL")
	configPath := "/config/appsettings.json"
	if isLocal == "true" {
		configPath = "appsettings.json"
	}

	if err := k.Load(file.Provider(configPath), json.Parser()); err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	var config config.Config
	if err := k.Unmarshal("", &config); err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}
	return config
}

// parseEventHandler returns the respective handler object specified in the config file.
func parseEventHandler(conf *config.Config) []handlers.Handler {

	var eventHandlers []handlers.Handler
	switch {
	case conf.Handler.Slack.Enabled && len(conf.Handler.Slack.Channel) > 0 || len(conf.Handler.Slack.Token) > 0:
		eventHandlers = append(eventHandlers, new(slack.Slack))
	case conf.Handler.SlackWebhook.Enabled && len(conf.Handler.SlackWebhook.Channel) > 0 || len(conf.Handler.SlackWebhook.Username) > 0 || len(conf.Handler.SlackWebhook.Slackwebhookurl) > 0:
		eventHandlers = append(eventHandlers, new(slackwebhook.SlackWebhook))
	case conf.Handler.Hipchat.Enabled && len(conf.Handler.Hipchat.Room) > 0 || len(conf.Handler.Hipchat.Token) > 0:
		eventHandlers = append(eventHandlers, new(hipchat.Hipchat))
	case conf.Handler.Mattermost.Enabled && len(conf.Handler.Mattermost.Channel) > 0 || len(conf.Handler.Mattermost.Url) > 0:
		eventHandlers = append(eventHandlers, new(mattermost.Mattermost))
	case conf.Handler.Flock.Enabled && len(conf.Handler.Flock.Url) > 0:
		eventHandlers = append(eventHandlers, new(flock.Flock))
	case conf.Handler.Webhook.Enabled && len(conf.Handler.Webhook.Url) > 0:
		eventHandlers = append(eventHandlers, new(webhook.Webhook))
	case conf.Handler.CloudEvent.Enabled && len(conf.Handler.CloudEvent.Url) > 0:
		eventHandlers = append(eventHandlers, new(cloudevent.CloudEvent))
	case conf.Handler.MSTeams.Enabled && len(conf.Handler.MSTeams.WebhookURL) > 0:
		eventHandlers = append(eventHandlers, new(msteam.MSTeams))
	case conf.Handler.SMTP.Enabled && len(conf.Handler.SMTP.Smarthost) > 0 || len(conf.Handler.SMTP.To) > 0:
		eventHandlers = append(eventHandlers, new(smtpClient.SMTP))
	case conf.Handler.Lark.Enabled && len(conf.Handler.Lark.WebhookURL) > 0:
		eventHandlers = append(eventHandlers, new(lark.Webhook))
	default:
		eventHandlers = append(eventHandlers, new(handlers.Default))
	}
	for _, eventHandler := range eventHandlers {
		if err := eventHandler.Init(conf); err != nil {
			logrus.Fatal(err)
		}
	}
	return eventHandlers
}
