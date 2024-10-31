package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildMessage(t *testing.T) {
	type args struct {
		template  string
		variables map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "multiple_variables",
			args: args{
				template: `hello $ONE, bye ${TWO}`,
				variables: map[string]string{
					"ONE": "tester1",
					"TWO": "tester2",
				},
			},
			want: `hello tester1, bye tester2`,
		},
		{
			name: "unknown_variable",
			args: args{
				template: `hello $UNKNOWN, bye ${SUSPECT}`,
				variables: map[string]string{
					"ONE": "tester1",
				},
			},
			want: `hello $UNKNOWN, bye ${SUSPECT}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildMessage(tt.args.template, tt.args.variables)
			assert.Equal(t, tt.want, got)
		})
	}
}
