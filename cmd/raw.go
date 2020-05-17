package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
		"request, key[:type]=value pairs. provide multiple times for multiple arguments. Without other options (see "+
		"--post, --delete), this will be a GET request.")
	rawCmd.Flags().BoolP("post", "p", false, "if true, REST request will be sent as a "+
		"POST. Cannot be used along with --websocket or --delete.")
	rawCmd.Flags().BoolP("delete", "d", false, "if true, REST request will be sent as a "+
		"DELETE. Cannot be used along with --websocket or --post.")
	Root.AddCommand(rawCmd)
}

type anyMessage struct {
	msgType string
	args    map[string]interface{}
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
	post, _ := cmd.Flags().GetBool("post")
	delete, _ := cmd.Flags().GetBool("delete")

	if ws && (post || delete) {
		logrus.Fatal("--post and --delete cannot be used along with --websocket")
	}

	if post && delete {
		logrus.Fatal("--post and --delete cannot be used together")
	}

	reqArgsList, _ := cmd.Flags().GetStringArray("arg")

	reqArgs := map[string]interface{}{}
	for _, arg := range reqArgsList {
		var key, typ, value string
		comps := strings.SplitN(arg, "=", 2)
		if len(comps) == 1 {
			logrus.Fatalf("expected key[:type]=value pair, but found %q", comps[0])
		}

		ktComps := strings.SplitN(comps[0], ":", 2)
		if len(ktComps) == 1 {
			ktComps = append(ktComps, "string")
		}
		key = ktComps[0]
		typ = ktComps[1]
		value = comps[1]

		switch typ {
		case "string":
			reqArgs[key] = value
		case "bool":
			b, err := strconv.ParseBool(value)
			if err != nil {
				logrus.WithError(err).Fatalf("%q: could not parse %q as bool", key, value)
			}
			reqArgs[key] = b
		default:
			logrus.Fatalf("unhandled type %q", typ)
		}
	}

	var result interface{}
	var err error
	if ws {
		msg := &anyMessage{msgType: args[0], args: reqArgs}
		result, err = client(cmd).RawWebsocketRequest(msg)
	} else {
		if post {
			result, err = client(cmd).RawRESTPost(args[0], reqArgs)
		} else if delete {
			result, err = client(cmd).RawRESTDelete(args[0], reqArgs)
		} else {
			result, err = client(cmd).RawRESTGet(args[0], reqArgs)
		}
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
