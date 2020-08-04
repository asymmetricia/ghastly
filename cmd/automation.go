package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pdbogen/ghastly/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

var automationCmd = &cobra.Command{
	Use:   "automation",
	Short: "sub-commands for manipulating and interacting with automations",
}

var automationListCmd = &cobra.Command{
	Use:   "list",
	Short: "retrieve a list of all known automations",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")
		ret, err := client(cmd).ListAutomations()
		if err != nil {
			logrus.WithError(err).Fatal("could not list config entries")
		}

		switch output {
		case "text":
			printTable(ret)
		case "json":
			retJson, _ := json.Marshal(ret)
			fmt.Println(string(retJson))
		}
	},
}

var automationGetCmd = &cobra.Command{
	Use:   "get [automation-id]",
	Short: "retrieve the configuration data for the given automation",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log := logrus.WithField("automation_id", args[0])
		ret, err := client(cmd).GetAutomation(api.AutomationId(args[0]))
		if err != nil {
			log.WithError(err).Fatal("could not get automation")
		}

		jb, err := json.Marshal(ret)
		if err != nil {
			log.WithError(err).Fatal("could not marshal automation to JSON")
		}
		fmt.Println(string(jb))
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		automations, err := client(cmd).ListAutomations()
		if err != nil {
			logrus.WithError(err).Fatal("could not get automations list")
			return nil, cobra.ShellCompDirectiveError
		}

		var ret []string
		for _, auto := range automations {
			if strings.HasPrefix(string(auto.Id), toComplete) {
				ret = append(ret, string(auto.Id))
			}
		}
		return ret, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	automationCmd.AddCommand(automationListCmd, automationGetCmd)
	Root.AddCommand(automationCmd)
}
