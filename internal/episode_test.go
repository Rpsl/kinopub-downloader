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
