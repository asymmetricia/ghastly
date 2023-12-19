package cmd

import (
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
		}
	},
}

func init() {
	entityCmd.AddCommand(entityGetCmd, entityListCmd)
	Root.AddCommand(entityCmd)
}
