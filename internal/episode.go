package internal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Episode struct {
	Title          string
	TVShow         string
	EpisodeNumber  int
	SeasonNumber   int
	URLForDownload string
	BasePath       string
}

const episodePattern = `s(\d+)e(\d+)`

func NewEpisode(title string, tvshow string, url string, basePath string) (*Episode, error) {
	ep := Episode{
		Title:          title,
		TVShow:         tvshow,
		URLForDownload: url,
		BasePath:       basePath,
	}

	num, err := ep.parseSeasonNumber(title)

	if err != nil {
		return nil, err
	}

	ep.SeasonNumber = num

	num, err = ep.parseEpisodeNumber(title)

	if err != nil {
		return nil, err
	}

	ep.EpisodeNumber = num

	return &ep, nil
}

func (e *Episode) GetPath() string {
	return fmt.Sprintf("%s/%s/Season %02d/%s.mp4", e.BasePath, e.pathEscape(e.TVShow), e.SeasonNumber, e.pathEscape(e.Title))
}

func (e *Episode) IsDownloaded() bool {
	_, err := os.Stat(e.GetPath())

	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func (e *Episode) GetURL() string {
	return e.URLForDownload
}

func (e *Episode) Download() (bool, error) {
	// todo need extract download engine into independent implementation
	err := e.makeSeasonDir()
	if err != nil {
		return false, err
	}

	out, err := os.Create(e.GetPath())
	if err != nil {
		return false, err
	}
	defer out.Close()

	resp, err := http.Get(e.URLForDownload)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)

	return true, err
}

func (e *Episode) pathEscape(path string) string {
	path = strings.ReplaceAll(path, "/", " ")

	reg := regexp.MustCompile("[\\s]+")
	path = reg.ReplaceAllString(path, " ")

	return path
}

func (e *Episode) makeSeasonDir() error {
	err := os.MkdirAll(fmt.Sprintf("%s/%s/Season %02d", e.BasePath, e.TVShow, e.SeasonNumber), os.ModePerm)
	if errors.Is(err, os.ErrExist) {
		return nil
	}

	return err
}

func (e *Episode) parseSeasonNumber(title string) (int, error) {
	r := regexp.MustCompile(episodePattern)

	match := r.FindAllStringSubmatch(title, 1)

	// todo it's look like shit
	if len(match) > 0 && len(match[0]) > 0 {
		n, err := strconv.Atoi(match[0][1])
		return n, err
	}

	return 0, errors.New(fmt.Sprintf("can't parse season number from title %s", title))
}

func (e *Episode) parseEpisodeNumber(title string) (int, error) {
	r := regexp.MustCompile(episodePattern)

	match := r.FindAllStringSubmatch(title, 1)

	// todo it's look like shit
	if len(match) > 0 && len(match[0]) > 1 {
		n, err := strconv.Atoi(match[0][2])
		return n, err
	}

	return 0, errors.New(fmt.Sprintf("can't parse episode number from title %s", title))
}
