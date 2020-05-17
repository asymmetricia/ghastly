package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/sirupsen/logrus"
)

// printTable prints the given table, which must be some kind of slice that can be unmarshalled to json and marshalled
// into []map[string]string. (for now, maybe clever reflection later)
func printTable(obj interface{}) {
	var table []map[string]string
	jsonBytes, err := json.Marshal(obj)
	if err == nil {
		err = json.Unmarshal(jsonBytes, &table)
	}
	if err != nil {
		logrus.Fatal(err)
	}

	columns := map[string]int{}
	for _, row := range table {
		for k, v := range row {
			if len(v) > columns[k] {
				columns[k] = len(v)
			}
		}
	}
	var sorted []string
	for k := range columns {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	for _, col := range sorted {
		fmt.Printf("%-"+strconv.Itoa(columns[col])+"s ", col)
	}

	fmt.Print("\n")
	for _, col := range sorted {
		for i := 0; i < columns[col]; i++ {
			fmt.Print("-")
		}
		fmt.Print(" ")
	}
	fmt.Print("\n")

	for _, row := range table {
		for _, col := range sorted {
			fmt.Printf("%-"+strconv.Itoa(columns[col])+"s ", row[col])
		}
		fmt.Print("\n")
	}
}
