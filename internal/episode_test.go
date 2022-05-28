package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEpisode(t *testing.T) {
	ep, err := NewEpisode("s01e03 - Blind Date", "30 Rock", "https://example.com/1.mp4", "/data")

	assert.Equal(t, 1, ep.SeasonNumber, "season number not equal")
	assert.Equal(t, 3, ep.EpisodeNumber, "episode number not equal")
	assert.Equal(t, ep.GetURL(), "https://example.com/1.mp4", "url of episode not equal")
	assert.Equal(t, ep.GetPath(), "/data/30 Rock/Season 01/s01e03 - Blind Date.mp4", "path for episode not equal")

	assert.Nil(t, err)
}

func TestNewEpisode_WithSlash(t *testing.T) {
	ep, err := NewEpisode("s03e01 - Where the Buffalo Roam/Venice Vengeance", "Close /Enough", "https://example.com/1.mp4", "/data")

	assert.Equal(t, 3, ep.SeasonNumber, "season number not equal")
	assert.Equal(t, 1, ep.EpisodeNumber, "episode number not equal")
	assert.Equal(t, ep.GetURL(), "https://example.com/1.mp4", "url of episode not equal")
	assert.Equal(t, ep.GetPath(), "/data/Close Enough/Season 03/s03e01 - Where the Buffalo Roam Venice Vengeance.mp4", "path for episode not equal")

	assert.Nil(t, err)
}

func TestNewEpisode_WithoutTitle(t *testing.T) {
	ep, err := NewEpisode("", "Close /Enough", "https://example.com/1.mp4", "/data")

	assert.Nil(t, ep, "episode object must be empty")
	assert.Error(t, err, "episode with empty title must return error")
}

func TestNewEpisode_WithoutShow(t *testing.T) {
	ep, err := NewEpisode("s03e01 - Where the Buffalo Roam/Venice Vengeance", "", "https://example.com/1.mp4", "/data")

	assert.Nil(t, ep, "episode object must be empty")
	assert.Error(t, err, "episode with empty show must return error")
}

func TestNewEpisode_WithoutUrl(t *testing.T) {
	ep, err := NewEpisode("s03e01 - Where the Buffalo Roam/Venice Vengeance", "Close /Enough", "", "/data")

	assert.Nil(t, ep, "episode object must be empty")
	assert.Error(t, err, "episode with empty url must return error")
}

func Test_parsePriority(t *testing.T) {
	var tests = []struct {
		name string
		args string
		want int
	}{
		{name: "normal", args: "s01e01 - Blind Date", want: 101},
		{name: "normal", args: "s01e03 - Blind Date", want: 103},
		{name: "normal", args: "s100e3 - Blind Date", want: 10003},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep, err := NewEpisode(tt.args, "30 Rock", "https://example.com/1.mp4", "/data")

			assert.Equalf(t, tt.want, ep.GetPriority(), "parsePriority(%v)")

			assert.Nil(t, err)
		})
	}
}
