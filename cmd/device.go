package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pdbogen/ghastly/search"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "sub-commands for manipulating and interacting with devices",
}

var deviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "retrieve a list of all known devices",
	Args:  cobra.NoArgs,
	Run:   runDeviceList,
}

var deviceSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "search for devices, using a simplified lucene syntax: `field:value {AND,OR} otherfield:othervalue`",
	Args:  cobra.MinimumNArgs(1),
	Run:   runDeviceSearch,
}

func runDeviceSearch(cmd *cobra.Command, args []string) {
	qs := strings.Join(args, " ")
	log := logrus.WithField("query", qs)
	queryI, err := search.Parse("query", []byte(qs))
	if err != nil {
		log.WithError(err).Fatal("could not parse query")
	}

	query, ok := queryI.(search.Node)
	if !ok {
		log.Fatalf("query parser returned nil error, but query was %T not search.Node", query)
	}

	devices, err := client(cmd).ListDevices()
	if err != nil {
		log.WithError(err).Fatal("could not get devices from HomeAssistant")
	}

	for _, device := range devices {
		attrs, err := search.Attributes(device)
		if err != nil {
			log.WithError(err).Fatalf("could not convert device %v to attributes: %v", device, err)
		}
		log.Debug(attrs)
		match, err := query.Evaluate(attrs)
		log.Debugf("result = %v, %v", match, err)
		if err != nil {
			log.WithError(err).Fatal("error during query evaluation")
		}
		if match {
			deviceJson, _ := json.Marshal(device)
			fmt.Println(string(deviceJson))
		}
	}
}

func runDeviceList(cmd *cobra.Command, args []string) {
	devices, err := client(cmd).ListDevices()
	if err != nil {
		logrus.WithError(err).Fatal("could not get devices from HomeAssistant")
	}
	devicesJson, _ := json.Marshal(devices)
	fmt.Println(string(devicesJson))
}

var deviceGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "retrieve a specific device with the given ID",
	Args:  cobra.ExactArgs(1),
	Run:   runDeviceGetCmd,
}

func runDeviceGetCmd(cmd *cobra.Command, args []string) {
	log := logrus.WithField("command", "device get")
	devices, err := client(cmd).ListDevices()
	if err != nil {
		log.Fatal(err)
	}
	if len(args) != 1 {
		log.Fatalf("expected one argument, got %d", len(args))
	}
	args[0] = strings.ToLower(args[0])
	for _, device := range devices {
		if strings.ToLower(device.ID) == args[0] {
			jsonBytes, err := json.Marshal(device)
			if err != nil {
				log.Fatalf("could not convert device to JSON: %v", err)
			}
			fmt.Println(string(jsonBytes))
			os.Exit(0)
		}
	}
}

func init() {
	deviceCmd.AddCommand(deviceGetCmd)
	deviceCmd.AddCommand(deviceListCmd)
	deviceCmd.AddCommand(deviceSearchCmd)
	Root.AddCommand(deviceCmd)
}

type Node interface{}
