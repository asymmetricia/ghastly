package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var entityCmd = &cobra.Command{
	Use:   "entity",
	Short: "sub-commands for manipulating and interacting with entities",
}

var entityGetCmd = &cobra.Command{
	Use:   "get [entity-id]",
	Short: "retrieve all known information about the given entity ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entity, err := client(cmd).GetEntity(args[0])
		if err != nil {
			logrus.Fatal(err)
		}
		fmt.Println(entity)
	},
}

var entityListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all known entities",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		entities, err := client(cmd).ListEntities()
		if err != nil {
			logrus.Fatal(err)
		}

		switch o, _ := cmd.Flags().GetString("output"); o {
		case "text":
			printTable(entities)
		case "json":
			jb, err := json.Marshal(entities)
			if err != nil {
				logrus.WithError(err).Fatal("could not marshal entity list to JSON")
			}
			fmt.Println(string(jb))
		}
	},
}

var entityRenameCmd = &cobra.Command{
	Use:   "rename [entity-id] [new-name]",
	Short: "rename an entity given by [entity-id] to have the friendly name given by [new-name]",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := client(cmd).SetEntityName(args[0], args[1]); err != nil {
			logrus.Fatal(err)
		}
	},
}

func init() {
	entityCmd.AddCommand(entityGetCmd, entityListCmd, entityRenameCmd)
	Root.AddCommand(entityCmd)
}
