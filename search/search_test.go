package search

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParse(t *testing.T) {
	for _, valid := range []string{
		"foo:bar",
		`"foo":"bar"`,
		`"\"foo\"":bar`,
		`"\"foo\"":"\\bar\\"`,
		"(foo:bar)",
		"((foo:bar))",
		"foo:bar AND bar:baz",
		"(foo:bar) AND (bar:baz)",
		"(foo:bar AND bar:baz) AND (blee:bloo)",
	} {
		_, err := Parse("test input", []byte(valid))
		require.NoError(t, err, "could not parse %q", valid)
		fmt.Println(valid, "ok")
	}
}
