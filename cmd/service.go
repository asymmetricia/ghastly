package cmd

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "sub-commands for manipulating and interacting with services",
}

var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "retrieves a list of services",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")
		ret, err := client(cmd).ListServices()
		if err != nil {
			logrus.WithError(err).Fatal("could not list services")
		}
		switch output {
		case "text":
			var table []map[string]string
			for _, svc := range ret {
				table = append(table, map[string]string{
					"domain": svc.Domain,
					"name":   svc.Name,
				})
			}
			sort.Slice(table, func(i, j int) bool {
				if table[i]["domain"] != table[j]["domain"] {
					return table[i]["domain"] < table[j]["domain"]
				}
				return table[i]["name"] < table[j]["name"]
			})
			printTable(table)
		case "json":
			retJson, err := json.Marshal(ret)
			if err != nil {
				logrus.WithError(err).Fatal("could not convert service list to JSON")
			}
			fmt.Println(string(retJson))
		}
	},
}

func init() {
	serviceCmd.AddCommand(
		serviceListCmd)
	Root.AddCommand(serviceCmd)
}
