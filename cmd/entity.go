package cmd

import (
	"fmt"
	"os"

	"github.com/pdbogen/ghastly/api"
	"github.com/spf13/cobra"
)

var entityCmd = &cobra.Command{
	Use:   "entity",
	Short: "sub-commands for manipulating and interacting with entities",
}

var entityGetCmd = &cobra.Command{
	Use:   "get",
	Short: "retrieve all known information about the given entity ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := &api.Client{Token: cmd.Flag("token").Value.String(), Server: cmd.Flag("server").Value.String()}
		entity, err := client.GetEntity(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(entity)
	},
}

func init() {
	entityCmd.AddCommand(entityGetCmd)
	Root.AddCommand(entityCmd)
}
