package status

import (
	"fmt"
	"strconv"
	"strings"
)

func Join(stringers ...fmt.Stringer) string {
	var stringBuilder strings.Builder
	for _, stringer := range stringers {
		stringBuilder.WriteString(stringer.String())
	}
	return stringBuilder.String()
}

type String string

func (s String) String() string { return string(s) }

type Variable string

func (v Variable) String() string {
	return fmt.Sprintf("#{%s}", string(v))
}

type Ternary struct {
	Variable string
	True     fmt.Stringer
	False    fmt.Stringer
}

func (t Ternary) String() string {
	return fmt.Sprintf("#{?%s,%s,%s}", t.Variable, t.True, t.False)
}

type Command struct {
	Name        string
	Args        []string
	Environment map[string]string
}

func (c Command) String() string {
	// Example: #(PATH="foo:$PATH" "$HOME/.dotfiles/bin/gowerline" status-right)
	var builder strings.Builder
	builder.WriteString("#(")
	for key, value := range c.Environment {
		builder.WriteString(fmt.Sprintf("%s=%q ", key, value))
	}
	builder.WriteString(strconv.Quote(c.Name))
	for _, arg := range c.Args {
		builder.WriteRune(' ')
		builder.WriteString(strconv.Quote(arg))
	}
	builder.WriteRune(')')
	return builder.String()
}
