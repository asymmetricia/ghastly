package cmd

import (
	"os"

	"github.com/pdbogen/ghastly/api"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Root = &cobra.Command{
	Use:   "ghastly",
	Short: "ghastly is a tool for interacting with homeassistant",
	Long: "A pretty incomplete tool for interacting with HomeAssistant. Mainly intended for exploring the API and " +
		"providing a test bed for the /api/ package. Future hopes includes developing a Terraform provider for " +
		"HomeAssistant.\n\nDownloads available on the GitHub Releases page: https://github.com/pdbogen/ghastly/releases",
}

var completion = &cobra.Command{
	Use:   "completion",
	Short: "generate bash completion scripts",
	Long: `To load completion, run:

. <(ghastly completion)

To configure your bash shell to load completions for each session add to your bashrc:

# ~/.bashrc or ~/.profile
. <(ghastly completion)
`,
	Run: func(cmd *cobra.Command, _ []string) {
		if err := Root.GenBashCompletion(os.Stdout); err != nil {
			logrus.WithError(err).Fatal("completion failed")
		}
	},
}

func client(cmd *cobra.Command) *api.Client {
	return &api.Client{Token: cmd.Flag("token").Value.String(), Server: cmd.Flag("server").Value.String()}
}

func init() {
	Root.AddCommand(completion)

	Root.PersistentFlags().String("token", os.Getenv("HASS_TOKEN"), "the bearer token used to authenticate to homeassistant. defaults to value of HASS_TOKEN environment variable")
	Root.PersistentFlags().String("server", os.Getenv("HASS_SERVER"), "the URL used to access homeassistant. defaults to value of HASS_SERVER environment variable")
	Root.PersistentFlags().String("loglevel", "INFO", "log level; one of TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC")
	Root.PersistentFlags().StringP("output", "o", "text", "output format for commands; options are `text` or `json`")

	cobra.OnInitialize(func() {
		lvlStr := Root.Flag("loglevel").Value.String()
		lvl, err := logrus.ParseLevel(lvlStr)
		if err != nil {
			logrus.Fatalf("bad log level %q: %v", lvlStr, err)
		}
		logrus.SetLevel(lvl)
	})
}
