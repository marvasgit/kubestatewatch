package client

import (
	"log"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/marvasgit/kubernetes-diffwatcher/config"
	"github.com/marvasgit/kubernetes-diffwatcher/pkg/controller"
	"github.com/marvasgit/kubernetes-diffwatcher/pkg/handlers/msteam"
)

func RunWithConfig() {

	conf := loadConfig()
	eventHandler := new(msteam.MSTeams)
	controller.Start(&conf, eventHandler)
}

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var k = koanf.New(".")

func loadConfig() config.Config {
	// Load JSON config.
	if err := k.Load(file.Provider("appsettings.json"), json.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	var config config.Config
	if err := k.Unmarshal("", &config); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	return config
}
