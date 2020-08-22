package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

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

		if len(args) == 0 {
			return completeServiceDomain(svcs)
		}

		r, d := completeServiceName(svcs, args[0])
		return r, d | cobra.ShellCompDirectiveNoSpace

	},
}

func completeServiceName(svcs []api.Service, domain string) ([]string, cobra.ShellCompDirective) {
	var ret []string
	for _, svc := range svcs {
		if domain == svc.Domain {
			ret = append(ret, svc.Name)
		}
	}
	sort.Strings(ret)
	return ret, cobra.ShellCompDirectiveNoFileComp
}

func completeServiceDomain(svcs []api.Service) ([]string, cobra.ShellCompDirective) {
	domains := map[string]bool{}
	for _, svc := range svcs {
		domains[svc.Domain] = true
	}

	var ret []string
	for domain := range domains {
		ret = append(ret, domain)
	}
	sort.Strings(ret)
	return ret, cobra.ShellCompDirectiveNoFileComp
}

var serviceCallCmd = &cobra.Command{
	Use: "call {domain} {service} [{field}={value} … {fieldN}={valueN}]",
	Short: "call the given service; each argument after service name should be " +
		"a `key=value` pair, which will be passed as a field in service_data",
	Args:              cobra.MinimumNArgs(2),
	Run:               serviceCallCmd_Run,
	ValidArgsFunction: serviceCallCmd_Complete,
}

func serviceCallCmd_Complete(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
	logrus.StandardLogger().SetOutput(os.Stderr)
	client := client(cmd)
	svcs, err := client.ListServices()
	if err != nil {
		logrus.WithError(err).Fatal("could not list services for tab completion")
	}

	if len(args) == 0 {
		return completeServiceDomain(svcs)
	}

	if len(args) == 1 {
		return completeServiceName(svcs, args[0])
	}

	var svc *api.Service
	for _, s := range svcs {
		if s.Domain == args[0] && s.Name == args[1] {
			svc = &s
			break
		}
	}

	if svc == nil {
		logrus.Errorf("no service %s.%s found", args[0], args[1])
		return nil, cobra.ShellCompDirectiveError |
			cobra.ShellCompDirectiveNoSpace |
			cobra.ShellCompDirectiveNoFileComp
	}

	var ret []string
	if !strings.Contains(complete, "=") {
		for k := range svc.Fields {
			ret = append(ret, k)
		}
		sort.Strings(ret)
		return ret, cobra.ShellCompDirectiveNoFileComp |
			cobra.ShellCompDirectiveNoSpace
	}

	fieldName := strings.SplitN(complete, "=", 2)[0]
	field, ok := svc.Fields[fieldName]
	if !ok {
		logrus.Errorf("service %s.%s has no field %q", args[0], args[1],
			fieldName)
		return nil, cobra.ShellCompDirectiveError |
			cobra.ShellCompDirectiveNoSpace |
			cobra.ShellCompDirectiveNoFileComp
	}
	
	if field.Type == api.Values {
		for _, v := range field.Values {
			ret = append(ret, fmt.Sprintf("%v", v))
		}
		return ret, cobra.ShellCompDirectiveNoFileComp
	}

	return nil, 0
}

func serviceCallCmd_Run(cmd *cobra.Command, args []string) {
	svc, err := client(cmd).GetService(args[0], args[1])
	if err != nil {
		logrus.WithError(err).Fatal("could not find service")
	}

	logrus.Tracef("call invoked with args %v", args)

	payload := map[string]interface{}{}
	for _, kv := range args[2:] {
		logrus.Tracef("parsing argument %q", kv)
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			logrus.Fatalf("%q was not `key=value` pair", kv)
		}

		field, ok := svc.Fields[parts[0]]
		if !ok {
			logrus.Fatal("service %s.%s does not have field %q", svc.Domain, svc.Name, parts[0])
		}
		logrus.Tracef("found field %v", field)

		switch field.Type {
		case api.Values:
			fallthrough
		case api.String:
			payload[parts[0]] = parts[1]
		case api.Number:
			payload[parts[0]], err = strconv.ParseFloat(parts[1], 64)
			if err != nil {
				logrus.WithError(err).Fatal("service %s.%s field %s is "+
					"numeric, but could not parse %q", svc.Domain, svc.Name,
					parts[0], parts[1])
			}
		case api.Boolean:
			payload[parts[0]], err = strconv.ParseBool(parts[1])
			if err != nil {
				logrus.WithError(err).Fatal("service %s.%s field %s is "+
					"boolean, but could not parse %q", svc.Domain, svc.Name,
					parts[0], parts[1])
			}
		default:
			logrus.Fatalf("field %q has unhandled type %s", field.Name, field.Type)
		}
	}

	logrus.Tracef("calling with payload %v", payload)
	states, err := svc.Call(payload)
	if err != nil {
		logrus.WithError(err).Fatalf("failed to call service %s.%s", svc.Domain, svc.Name)
	}
	statesJson, err := json.Marshal(states)
	if err != nil {
		logrus.WithError(err).Fatal("could not marshal states to JSON: %w", err)
	}
	fmt.Println(string(statesJson))
}

func init() {
	serviceCmd.AddCommand(
		serviceListCmd,
		serviceGetCmd,
		serviceCallCmd)
	Root.AddCommand(serviceCmd)
}
