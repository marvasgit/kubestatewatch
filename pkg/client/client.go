package client

import (
	"github.com/marvasgit/kubernetes-diffwatcher/config"
	"github.com/marvasgit/kubernetes-diffwatcher/pkg/controller"
	"github.com/marvasgit/kubernetes-diffwatcher/pkg/handlers/msteam"
)

func RunWithConfig(conf *config.Config) {

	eventHandler := new(msteam.MSTeams)
	//TODO: fix it without handler, and extract the handling to have abiliy to push to more then one place
	controller.Start(conf, eventHandler)
}
