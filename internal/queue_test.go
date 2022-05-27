package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parsePriority1(t *testing.T) {
	var tests = []struct {
		name string
		args string
		want int64
	}{
		{name: "normal", args: "s01e01 - Blind Date", want: 101},
		{name: "normal", args: "s01e03 - Blind Date", want: 103},
		{name: "normal", args: "s100e3 - Blind Date", want: 10003},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep, err := NewEpisode(tt.args, "30 Rock", "https://example.com/1.mp4", "/data")

			assert.Equalf(t, tt.want, parsePriority(ep), "parsePriority(%v)")

			assert.Nil(t, err)
		})
	}
}
