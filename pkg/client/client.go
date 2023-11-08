package client

import (
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/marvasgit/kubernetes-diffwatcher/config"
	"github.com/marvasgit/kubernetes-diffwatcher/pkg/controller"
	"github.com/marvasgit/kubernetes-diffwatcher/pkg/handlers"
	"github.com/marvasgit/kubernetes-diffwatcher/pkg/handlers/msteam"
	"github.com/sirupsen/logrus"
)

func Run() {

	conf := loadConfig()
	eventHandler := new(msteam.MSTeams)
	//TODO: fix it without handler, and extract the handling to have abiliy to push to more then one place
	controller.Start(&conf, []handlers.Handler{eventHandler})
}

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var k = koanf.New(".")

func loadConfig() config.Config {
	// Load JSON config.
	if err := k.Load(file.Provider("appsettings.json"), json.Parser()); err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	var config config.Config
	if err := k.Unmarshal("", &config); err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}
	return config
}
