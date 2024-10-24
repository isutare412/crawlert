package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_escapeForMarkdown(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "sample_messages",
			args: args{
				s: `Hello, world!`,
			},
			want: `Hello, world\!`,
		},
		{
			name: "special_characters",
			args: args{
				s: `*[]()~>#+-=|{}.!`,
			},
			want: `\*\[\]\(\)\~\>\#\+\-\=\|\{\}\.\!`,
		},
		{
			name: "message_with_link",
			args: args{
				s: `This is Google: https://google.com`,
			},
			want: `This is Google: https://google\.com`,
		},
		{
			name: "some_characters_are_allowed",
			args: args{
				s: "__Hello__, _world_, `block`",
			},
			want: "__Hello__, _world_, `block`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeForMarkdown(tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
