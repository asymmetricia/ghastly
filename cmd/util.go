package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	prettyTable "github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/sirupsen/logrus"
)

type tableCell string

func (t *tableCell) UnmarshalJSON(data []byte) error {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	switch o := obj.(type) {
	case string:
		*t = tableCell(o)
	case float64:
		*t = tableCell(fmt.Sprintf("%f", o))
	case []interface{}:
		*t = tableCell(data)
	case nil:
	default:
		return fmt.Errorf("cannot handle json-decoded object %v of type %T", obj, obj)
	}

	return nil
}

var _ json.Unmarshaler = (*tableCell)(nil)

// printTable prints the given table, which must be some kind of slice that can
// be unmarshalled to json and marshalled into []map[string]string. (for now,
// maybe clever reflection later) If `order` is provided, columns will be
// ordered as represented in the slice. Any additional columns will be ordered
// arbitrarily at the end. As a special case, any column named `description`
// will be wrapped at 40 columns.
func printTable(obj interface{}, order ...[]string) {
	if len(order) == 0 {
		order = [][]string{{}}
	}

	var table []map[string]tableCell
	jsonBytes, err := json.Marshal(obj)
	if err == nil {
		err = json.Unmarshal(jsonBytes, &table)
	}
	if err != nil {
		logrus.Fatalf("unmarshaling to convert to table: %v", err)
	}
	if len(table) == 0 {
		return
	}

	ordered := map[string]bool{}
	var sorted []string
	for _, column := range order[0] {
		sorted = append(sorted, column)
		ordered[column] = true
	}
	for k := range table[0] {
		if !ordered[k] {
			sorted = append(sorted, k)
		}
	}
	sort.Strings(sorted[len(ordered):])

	t := prettyTable.NewWriter()
	t.SetOutputMirror(os.Stdout)

	hrow := make([]interface{}, len(sorted))
	for i, v := range sorted {
		hrow[i] = v
	}
	t.AppendHeader(hrow)

	for _, row := range table {
		var pRow prettyTable.Row
		for _, column := range sorted {
			pRow = append(pRow, text.WrapText(
				text.WrapSoft(
					string(row[column]),
					40),
				40))
		}
		t.AppendRow(pRow)
	}
	t.Render()
}
