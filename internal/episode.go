package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Episode struct {
	Title          string
	Show           string
	EpisodeNumber  int
	SeasonNumber   int
	URLForDownload string
	BasePath       string
}

const episodePattern = `s(\d+)e(\d+)`

func NewEpisode(title string, show string, url string, basePath string) (*Episode, error) {
	switch {
	case title == "":
		return nil, errors.New("can't process episode without title")
	case show == "":
		return nil, errors.New("can't process episode without show title")
	case url == "":
		return nil, errors.New("can't process episode without url")
	}

	ep := Episode{
		Title:          title,
		Show:           show,
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
	return fmt.Sprintf("%s/%s/Season %02d/%s.mp4", e.BasePath, e.pathEscape(e.Show), e.SeasonNumber, e.pathEscape(e.Title))
}

func (e *Episode) IsDownloaded() bool {
	_, err := os.Stat(e.GetPath())

	return !errors.Is(err, os.ErrNotExist)
}

func (e *Episode) GetURL() string {
	return e.URLForDownload
}

func (e *Episode) Download(ctx context.Context) (bool, error) {
	log.Infof("[+++] start downloading - \"%s - %s\" into \"%s\"", e.TVShow, e.Title, e.GetPath())

	// todo need extract download engine into independent implementation
	err := e.makeSeasonDir()
	if err != nil {
		log.WithError(err).Errorf("can't create season dir: \"%s\"", e.GetPath())
		return false, err
	}

	file, err := os.Create(e.GetPath())
	if err != nil {
		log.WithError(err).Errorf("can't create episode file: \"%s\"", e.GetPath())
		return false, err
	}
	defer file.Close()

	// todo check file length
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, e.URLForDownload, nil)
	resp, urlErr := http.DefaultClient.Do(req)

	if urlErr != nil {
		e.removeTempFile(file)
		log.WithError(urlErr).Errorf("error while doing request - \"%s - %s\"", e.TVShow, e.Title)

		return false, urlErr
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		e.removeTempFile(file)
		log.WithError(err).Errorf("error while downloading - \"%s - %s\"", e.TVShow, e.Title)
	}

	return true, err
}

func (e *Episode) removeTempFile(file *os.File) {
	// remove incomplete file when caught error
	file.Close()
	os.Remove(e.GetPath())
}

func (e *Episode) pathEscape(path string) string {
	path = strings.ReplaceAll(path, "/", " ")

	reg := regexp.MustCompile(`\s+`)
	path = reg.ReplaceAllString(path, " ")

	return path
}

func (e *Episode) makeSeasonDir() error {
	err := os.MkdirAll(fmt.Sprintf("%s/%s/Season %02d", e.BasePath, e.Show, e.SeasonNumber), os.ModePerm)
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

	return 0, fmt.Errorf("can't parse season number from title \"%s\"", title)
}

func (e *Episode) parseEpisodeNumber(title string) (int, error) {
	r := regexp.MustCompile(episodePattern)

	match := r.FindAllStringSubmatch(title, 1)

	// todo it's look like shit
	if len(match) > 0 && len(match[0]) > 1 {
		n, err := strconv.Atoi(match[0][2])
		return n, err
	}

	return 0, fmt.Errorf("can't parse episode number from title \"%s\"", title)
}
