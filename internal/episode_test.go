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
