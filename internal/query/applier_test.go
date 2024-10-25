package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/isutare412/crawlert/internal/core/model"
)

var rawJSONs = [...]string{
	`[
  {
    "name": "apple",
    "color": "green",
    "price": 1.2
  },
  {
    "name": "banana",
    "color": "yellow",
    "price": 0.5
  },
  {
    "name": "kiwi",
    "color": "green",
    "price": 1.25
  }
]`,
	`{
  "people": [
    {
      "name": "Alice",
      "age": 20,
      "friends": [
        {
          "name": "friend-one",
          "relationship": "good"
        }
      ]
    },
    {
      "name": "Bob",
      "age": 24,
      "friends": [
        {
          "name": "friend-two",
          "relationship": "poor"
        },
        {
          "name": "friend-three",
          "relationship": "bad"
        }
      ]
    }
  ]
}`,
}

func TestExecutor_ApplyQuery(t *testing.T) {
	type inits struct {
		checkQuery      string
		variableQueries map[string]string
	}
	type args struct {
		jsonBytes []byte
	}
	tests := []struct {
		name    string
		inits   inits
		args    args
		want    model.QueryResult
		wantErr bool
	}{
		{
			name: "select_length",
			inits: inits{
				checkQuery: `[ .[] | select(.name == "apple") ] | length`,
			},
			args: args{
				jsonBytes: []byte(rawJSONs[0]),
			},
			want: model.QueryResult{
				Matched:   true,
				Variables: map[string]string{},
			},
		},
		{
			name: "not_found_by_select",
			inits: inits{
				checkQuery: `[ .[] | select(.foo == "bar") ] | length`,
			},
			args: args{
				jsonBytes: []byte(rawJSONs[0]),
			},
			want: model.QueryResult{
				Matched:   false,
				Variables: map[string]string{},
			},
		},
		{
			name: "variables_are_set",
			inits: inits{
				checkQuery: `.people[] | select( .friends | length >= 2 ) | length`,
				variableQueries: map[string]string{
					"NAMES":       `[ .people[] | .name ]`,
					"BAD_FREINDS": `[ .people[] | { "name": .name, "badFriends": [ .friends[] | select( .relationship == "bad" or .relationship == "poor" ) ] } ]`,
				},
			},
			args: args{
				jsonBytes: []byte(rawJSONs[1]),
			},
			want: model.QueryResult{
				Matched: true,
				Variables: map[string]string{
					"NAMES":       `["Alice","Bob"]`,
					"BAD_FREINDS": `[{"badFriends":[],"name":"Alice"},{"badFriends":[{"name":"friend-two","relationship":"poor"},{"name":"friend-three","relationship":"bad"}],"name":"Bob"}]`,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewApplier(tt.inits.checkQuery, tt.inits.variableQueries)
			require.NoError(t, err)

			resp, err := e.ApplyQuery(tt.args.jsonBytes)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, resp)
		})
	}
}

func Test_isTruthyValue(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "boolean_true",
			args: args{
				s: "true",
			},
			want: true,
		},
		{
			name: "boolean_false",
			args: args{
				s: "false",
			},
			want: false,
		},
		{
			name: "positive_number",
			args: args{
				s: "412",
			},
			want: true,
		},
		{
			name: "negative_number",
			args: args{
				s: "-412",
			},
			want: false,
		},
		{
			name: "zero",
			args: args{
				s: "0",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTruthyValue(tt.args.s); got != tt.want {
				t.Errorf("isTruthyValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
