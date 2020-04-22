package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pdbogen/ghastly/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "sub-commands for manipulating and interacting with config entries",
}

var configListEntriesCmd = &cobra.Command{
	Use:   "list-entries",
	Short: "retrieve a list of all known config entries",
	Args:  cobra.NoArgs,
	Run:   runConfigListEntries,
}

var configListFlowsCmd = &cobra.Command{
	Use:   "list-flows",
	Short: "retrieve a list of all known config flows",
	Args:  cobra.NoArgs,
	Run:   runConfigListFlows,
}

func runConfigListFlows(cmd *cobra.Command, args []string) {
	client := &api.Client{Token: cmd.Flag("token").Value.String(), Server: cmd.Flag("server").Value.String()}
	ret, err := client.ListConfigFlows()
	if err != nil {
		logrus.WithError(err).Fatal("could not list config flows")
	}
	retJson, _ := json.Marshal(ret)
	fmt.Println(string(retJson))
}

func init() {
	configCmd.AddCommand(configListEntriesCmd)
	Root.AddCommand(configCmd)
}

func runConfigListEntries(cmd *cobra.Command, args []string) {
	client := &api.Client{Token: cmd.Flag("token").Value.String(), Server: cmd.Flag("server").Value.String()}
	ret, err := client.ListConfigEntries()
	if err != nil {
		logrus.WithError(err).Fatal("could not list config entries")
	}
	retJson, _ := json.Marshal(ret)
	fmt.Println(string(retJson))
}
