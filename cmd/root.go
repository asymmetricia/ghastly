package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Root = &cobra.Command{
	Use:   "ghastly",
	Short: "ghastly is a tool for interacting with homeassistant",
}

func init() {
	Root.PersistentFlags().String("token", os.Getenv("HASS_TOKEN"), "the bearer token used to authenticate to homeassistant. defaults to value of HASS_TOKEN environment variable")
	Root.PersistentFlags().String("server", os.Getenv("HASS_SERVER"), "the URL used to access homeassistant. defaults to value of HASS_SERVER environment variable")
	Root.PersistentFlags().String("loglevel", "INFO", "log level; one of TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC")

	cobra.OnInitialize(func() {
		lvlStr := Root.Flag("loglevel").Value.String()
		lvl, err := logrus.ParseLevel(lvlStr)
		if err != nil {
			logrus.Fatalf("bad log level %q: %v", lvlStr, err)
		}
		logrus.SetLevel(lvl)
	})
}
