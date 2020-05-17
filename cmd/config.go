package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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
	Run: func(cmd *cobra.Command, args []string) {
		ret, err := client(cmd).ListConfigEntries()
		if err != nil {
			logrus.WithError(err).Fatal("could not list config entries")
		}
		retJson, _ := json.Marshal(ret)
		fmt.Println(string(retJson))
	},
}

var configListFlowsCmd = &cobra.Command{
	Use:   "list-flows",
	Short: "list in-progress but not started config flows",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ret, err := client(cmd).ListConfigFlowProgress()
		if err != nil {
			logrus.WithError(err).Fatal("could not list config flows")
		}
		retJson, _ := json.Marshal(ret)
		fmt.Println(string(retJson))
	},
}

var configGetFlowCmd = &cobra.Command{
	Use:   "get-flow [id]",
	Short: "retrieve the config flow associated with the given ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ret, err := client(cmd).GetFlow(api.ConfigFlowId(args[0]))
		if err != nil {
			logrus.WithError(err).Fatalf("could not get config flow %q", args[0])
		}
		retJson, _ := json.Marshal(ret)
		fmt.Println(string(retJson))
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) (i []string, directive cobra.ShellCompDirective) {
		flows, err := client(cmd).ListConfigFlowProgress()
		if err != nil {
			logrus.WithError(err).Error("could not list flows")
			return nil, cobra.ShellCompDirectiveError | cobra.ShellCompDirectiveNoFileComp
		}
		var ret []string
		for _, flow := range flows {
			if strings.HasPrefix(string(flow.FlowId), toComplete) {
				ret = append(ret, string(flow.FlowId))
			}
		}
		return ret, cobra.ShellCompDirectiveNoFileComp
	},
}

var configSetFlowCmd = &cobra.Command{
	Use:   "set-flow [flow-id] [key=value] ... [key=value]",
	Short: "set the given key-value pairs on the flow with the given ID",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log := logrus.WithField("flow_id", args[0])
		id := api.ConfigFlowId(args[0])
		kv := map[string]string{}
		for _, pair := range args[1:] {
			comps := strings.SplitN(pair, "=", 2)
			if len(comps) != 2 {
				log.Fatalf("expected key=value pair but %q did not contain `=`", pair)
			}
			kv[comps[0]] = comps[1]
		}

		payload := map[string]interface{}{}

		flow, err := client(cmd).GetFlow(id)
		if err != nil {
			log.WithError(err).Fatal("could not get flow")
		}
		for _, field := range flow.DataSchema {
			if _, ok := kv[field.Name]; !ok {
				if field.Required {
					log.Fatalf("flow requires field %q but was missing from set command", field.Name)
				}
				continue
			}
			switch field.Type {
			case "string":
				payload[field.Name] = kv[field.Name]
			case "integer":
				i, err := strconv.Atoi(kv[field.Name])
				if err != nil {
					log.WithError(err).Fatalf("field %q expected integer, but %q could not be parsed as integer",
						field.Name, kv[field.Name])
				}
				payload[field.Name] = i
			case "boolean":
				b, err := strconv.ParseBool(kv[field.Name])
				if err != nil {
					log.WithError(err).Fatalf("Field %q expected boolean, but %q could not be parsed as bool",
						field.Name, kv[field.Name])
				}
				payload[field.Name] = b
			default:
				log.Fatalf("not sure how to handle field type %q", field.Type)
			}
		}

		log.Debugf("prepared payload: %v", payload)

		result, err := client(cmd).SetFlow(id, payload)
		if err != nil {
			log.WithError(err).Fatal("set-flow failed")
		}
		log.WithField("result", result).Print("set successful")
	},
}

var configStartOptionsFlowCmd = &cobra.Command{
	Use:   "start-options-flow [handler-id]",
	Short: "start a config flow with with the given handler ID as the handler",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ret, err := client(cmd).StartOptionsFlow(api.EntryId(args[0]))
		if err != nil {
			logrus.WithError(err).Fatal("start-flow failed")
		}
		retJson, _ := json.Marshal(ret)
		fmt.Println(string(retJson))
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		entries, err := client(cmd).ListConfigEntries()
		if err != nil {
			logrus.WithError(err).Error("could not list config entries")
			return nil, cobra.ShellCompDirectiveError
		}
		var ret []string
		for _, entry := range entries {
			if strings.HasPrefix(string(entry.EntryId), toComplete) {
				ret = append(ret, string(entry.EntryId))
			}
		}
		return ret, cobra.ShellCompDirectiveNoFileComp
	},
}

var configGetOptionsFlowCmd = &cobra.Command{
	Use:   "get-options-flow [flow-id]",
	Short: "retrieve the current state of the indicated options flow",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ret, err := client(cmd).GetOptionsFlow(api.ConfigFlowId(args[0]))
		if err != nil {
			logrus.WithError(err).Fatal("get-options-flow failed")
		}
		retJson, _ := json.Marshal(ret)
		fmt.Println(string(retJson))
	},
}

var configGetConfigCmd = &cobra.Command{
	Use:   "get-config",
	Short: "retrieve the top-level system config",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		ret, err := client(cmd).GetConfig()
		if err != nil {
			logrus.WithError(err).Fatal("could not get config")
		}
		retJson, _ := json.Marshal(ret)
		fmt.Println(string(retJson))
	},
}

var configListFlowHandlersCmd = &cobra.Command{
	Use:   "list-flow-handlers",
	Short: "returns a list of flow handlers, whatever those are",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		ret, err := client(cmd).ListFlowHandlers()
		if err != nil {
			logrus.WithError(err).Fatal("could not list flow handlers")
		}
		retJson, _ := json.Marshal(ret)
		fmt.Println(string(retJson))
	},
}

func init() {
	configCmd.AddCommand(configListEntriesCmd,
		configListFlowHandlersCmd,
		configListFlowsCmd,
		configGetFlowCmd,
		configGetConfigCmd,
		configGetOptionsFlowCmd,
		configSetFlowCmd,
		configStartOptionsFlowCmd)
	Root.AddCommand(configCmd)
}
