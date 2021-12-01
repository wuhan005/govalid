package govalid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseRules(t *testing.T) {
	for _, tc := range []struct {
		name string
		rule string
		want []*rule
	}{
		{
			name: "single rule",
			rule: "max:5",
			want: []*rule{
				{checker: "max", params: []string{"5"}},
			},
		},
		{
			name: "multiple rule",
			rule: "required;max:5;min:0",
			want: []*rule{
				{checker: "required"},
				{checker: "max", params: []string{"5"}},
				{checker: "min", params: []string{"0"}},
			},
		},
		{
			name: "no value",
			rule: "required:",
			want: []*rule{
				{checker: "required", params: []string{""}},
			},
		},
		{
			name: "nothing",
			rule: "::::",
			want: []*rule{},
		},
		{
			name: "empty",
			rule: "",
			want: []*rule{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := parseRules(tc.rule)
			assert.Equal(t, tc.want, got)
		})
	}
}
