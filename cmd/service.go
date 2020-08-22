package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/pdbogen/ghastly/api"
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

var serviceGetCmd = &cobra.Command{
	Use:   "get {domain} {service-name}",
	Short: "retrieve the given service in the given domain",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")
		ret, err := client(cmd).GetService(args[0], args[1])
		if err != nil {
			logrus.WithError(err).Fatal("could not get service")
		}
		switch output {
		case "text":
			fmt.Println("Domain:", ret.Domain)
			fmt.Println("Name:  ", ret.Name)
			if len(ret.Fields) == 0 {
				fmt.Println("No fields.")
			} else {
				fmt.Println("Fields:")
				var fields []*api.ServiceField
				for _, field := range ret.Fields {
					fields = append(fields, field)
				}
				sort.Slice(fields, func(i, j int) bool {
					return fields[i].Name < fields[j].Name
				})
				printTable(fields, []string{"name", "description"})
			}
		case "json":
			retJson, err := json.Marshal(ret)
			if err != nil {
				logrus.WithError(err).Fatal("could not convert service list to JSON")
			}
			fmt.Println(string(retJson))
		}
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		logrus.StandardLogger().SetOutput(os.Stderr)
		svcs, err := client(cmd).ListServices()
		if err != nil {
			logrus.WithError(err).Fatal("could not list services for tab completion")
		}

		if len(args) > 1 {
			return nil, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
		}

		var ret []string
		if len(args) == 0 {
			domains := map[string]bool{}
			for _, svc := range svcs {
				domains[svc.Domain] = true
			}
			for domain := range domains {
				ret = append(ret, domain)
			}
			sort.Strings(ret)
			return ret, cobra.ShellCompDirectiveNoFileComp
		}

		for _, svc := range svcs {
			if args[0] == svc.Domain {
				ret = append(ret, svc.Name)
			}
		}
		sort.Strings(ret)
		return ret, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
	},
}

func init() {
	serviceCmd.AddCommand(
		serviceListCmd,
		serviceGetCmd)
	Root.AddCommand(serviceCmd)
}
