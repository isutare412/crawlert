package discord

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_splitMessage(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want []string
	}{
		{
			name: "empty_message",
			msg:  "",
			want: []string{""},
		},
		{
			name: "short_message",
			msg:  "hello world",
			want: []string{"hello world"},
		},
		{
			name: "exactly_max_length",
			msg:  strings.Repeat("a", maxMessageLength),
			want: []string{strings.Repeat("a", maxMessageLength)},
		},
		{
			name: "exceeds_max_splits_at_newline",
			msg:  strings.Repeat("a", maxMessageLength-10) + "\n" + strings.Repeat("b", 20),
			want: []string{
				strings.Repeat("a", maxMessageLength-10) + "\n",
				strings.Repeat("b", 20),
			},
		},
		{
			name: "exceeds_max_no_newline",
			msg:  strings.Repeat("a", maxMessageLength+500),
			want: []string{
				strings.Repeat("a", maxMessageLength),
				strings.Repeat("a", 500),
			},
		},
		{
			name: "splits_into_three_chunks",
			msg:  strings.Repeat("a", maxMessageLength-1) + "\n" + strings.Repeat("b", maxMessageLength-1) + "\n" + strings.Repeat("c", 100),
			want: []string{
				strings.Repeat("a", maxMessageLength-1) + "\n",
				strings.Repeat("b", maxMessageLength-1) + "\n",
				strings.Repeat("c", 100),
			},
		},
		{
			name: "prefers_last_newline_in_chunk",
			msg:  strings.Repeat("a", maxMessageLength-10) + "\n" + strings.Repeat("b", 5) + "\n" + strings.Repeat("c", 100),
			want: []string{
				strings.Repeat("a", maxMessageLength-10) + "\n" + strings.Repeat("b", 5) + "\n",
				strings.Repeat("c", 100),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitMessage(tt.msg)
			assert.Equal(t, tt.want, got)

			for i, chunk := range got {
				assert.LessOrEqual(t, len(chunk), maxMessageLength, "chunk %d exceeds max length", i)
			}

			assert.Equal(t, tt.msg, strings.Join(got, ""), "joined chunks should equal original message")
		})
	}
}
