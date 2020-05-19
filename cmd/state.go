package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "sub-commands for manipulating and interacting with device states",
}

var stateListCmd = &cobra.Command{
	Use:   "list",
	Short: "retrieve a list of all known device states",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ret, err := client(cmd).ListStates()
		if err != nil {
			logrus.WithError(err).Fatal("could not list states")
		}
		retJson, _ := json.Marshal(ret)
		fmt.Println(string(retJson))
	},
}

func init() {
	stateCmd.AddCommand(stateListCmd)
	Root.AddCommand(stateCmd)
}
