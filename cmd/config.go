package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pdbogen/ghastly/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
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
	Short: "retrieve a list of all known config flows",
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
					log.WithError(err).Fatal("field %q expected integer, but %q could not be parsed as integer",
						field.Name, kv[field.Name])
				}
				payload[field.Name] = i
			case "boolean":
				b, err := strconv.ParseBool(kv[field.Name])
				if err != nil {
					log.WithError(err).Fatal("Field %q expected boolean, but %q could not be parsed as bool",
						field.Name, kv[field.Name])
				}
				payload[field.Name] = b
			default:
				log.Fatalf("not sure how to handle field type %q", field.Type)
			}
		}

		log.Debugf("prepared payload: %v", payload)

		if err := client(cmd).SetFlow(id, payload); err != nil {
			log.WithError(err).Fatal("set-flow failed")
		}
		log.Print("set successful")
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

func init() {
	configCmd.AddCommand(configListEntriesCmd,
		configListFlowsCmd,
		configGetFlowCmd,
		configGetConfigCmd,
		configSetFlowCmd)
	Root.AddCommand(configCmd)
}
