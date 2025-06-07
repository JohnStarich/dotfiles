package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		description string
		in          string
		args        []string
		out         string
		errOut      string
	}{
		{
			description: "one valid line",
			in:          `{"foo":"bar"}`,
			out: `{
  "foo": "bar"
}
`,
		},
		{
			description: "one valid and one invalid line",
			in: `
foo
{"foo":"bar"}
`,
			out: `{
  "_json_parse_error": "foo"
}
{
  "foo": "bar"
}
`,
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			t.Cleanup(cancel)
			var output, errOutput bytes.Buffer
			err := run(ctx, tc.args, strings.NewReader(tc.in), &output, &errOutput)
			assert.NoError(t, err)
			assert.Equal(t, tc.out, output.String())
			assert.Equal(t, tc.errOut, errOutput.String())
		})
	}
}
