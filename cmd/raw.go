package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

var rawCmd = &cobra.Command{
	Use:   "raw <path>",
	Short: "send a raw request and print the result",
	Args:  cobra.ExactArgs(1),
	Run:   runRaw,
}

func init() {
	rawCmd.Flags().BoolP("websocket", "w", false, "if true, send request to the given "+
		"websocket endpoint; otherwise, send a GET")
	rawCmd.Flags().StringArrayP("arg", "a", nil, "arguments to send along with the "+
		"request, key=value pairs. provide multiple times for multiple argments.")
	Root.AddCommand(rawCmd)
}

type anyMessage struct {
	msgType string
	args    map[string]string
}

func (a *anyMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.args)
}

var _ json.Marshaler = (*anyMessage)(nil)

func (a *anyMessage) Type() string {
	return a.msgType
}

func runRaw(cmd *cobra.Command, args []string) {
	ws, _ := cmd.Flags().GetBool("websocket")
	reqArgsList, _ := cmd.Flags().GetStringArray("arg")

	reqArgs := map[string]string{}
	for _, arg := range reqArgsList {
		comps := strings.Split(arg, "=")
		reqArgs[comps[0]] = strings.Join(comps[1:], "=")
	}

	var result interface{}
	var err error
	if ws {
		msg := &anyMessage{msgType: args[0], args: reqArgs}
		result, err = client(cmd).RawWebsocketRequest(msg)
	} else {
		result, err = client(cmd).RawRESTRequest(args[0], reqArgs)
	}
	if err != nil {
		logrus.WithError(err).Fatal("could not send request")
	}

	rj, err := json.Marshal(result)
	if err != nil {
		logrus.WithError(err).Fatal("could not marshal result to JSON")
	}
	fmt.Println(string(rj))
}
