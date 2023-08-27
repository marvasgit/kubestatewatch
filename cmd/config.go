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

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/marvasgit/diffwatcher/config"
	"github.com/marvasgit/diffwatcher/pkg/client"
	"github.com/marvasgit/diffwatcher/pkg/event"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const diffwatcherConfigFile = ".diffwatcher.yaml"

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "modify diffwatcher configuration",
	Long: `
config command allows configuration of ~/.diffwatcher.yaml for running diffwatcher`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var configAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add webhook config to ~/.diffwatcher.yaml",
	Long: `
Adds webhook config to ~/.diffwatcher.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var configTestCmd = &cobra.Command{
	Use:   "test",
	Short: "test handler config present in ~/.diffwatcher.yaml",
	Long: `
Tests handler configs present in ~/.diffwatcher.yaml by sending test messages`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Testing Handler configs from .diffwatcher.yaml")
		conf, err := config.New()
		if err != nil {
			logrus.Fatal(err)
		}
		eventHandler := client.ParseEventHandler(conf)
		e := event.DiffWatchEvent{
			Namespace: "testNamespace",
			Name:      "testResource",
			Kind:      "testKind",
			Component: "testComponent",
			Host:      "testHost",
			Reason:    "Tested",
			Status:    "Normal",
		}
		eventHandler.Handle(e)
	},
}

var configSampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "Show a sample config file",
	Long: `
Print a sample config file which can be put in ~/.diffwatcher.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(config.ConfigSample)
	},
}

var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "view ~/.diffwatcher.yaml",
	Long: `
Display the contents of the contents of ~/.diffwatcher.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(os.Stderr, "Contents of ~/.diffwatcher.yaml")
		configFile, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), diffwatcherConfigFile))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(configFile))
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(
		configAddCmd,
		configTestCmd,
		configSampleCmd,
		configViewCmd,
	)

	configAddCmd.AddCommand(
		slackConfigCmd,
		slackwebhookConfigCmd,
		hipchatConfigCmd,
		mattermostConfigCmd,
		flockConfigCmd,
		webhookConfigCmd,
		cloudEventConfigCmd,
		msteamsConfigCmd,
		smtpConfigCmd,
		larkConfigCmd,
	)
}
