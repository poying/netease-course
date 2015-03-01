package necourse

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

var (
	playListIdRe = regexp.MustCompile(`data-plid='([^']+)'`)
	videoIdRe    = regexp.MustCompile(`/([-\w]+)_([-\w]+)\.html$`)
	listPageRe   = regexp.MustCompile(`/special/(?:[-\w]+)/(?:[-\w]+)\.html$`)
	videoPageRe  = regexp.MustCompile(`/movie/\d{4}/\d{1,2}/[^/]/[^/]/[^/]+\.html$`)
)

func Get(url string) (Course, error) {
	if isVideoPage(url) {
		return getVideosFromVideoPage(url)
	}

	if isListPage(url) {
		return getVideosFromListPage(url)
	}

	return nil, errors.New("Unknown page: " + url)
}

func isVideoPage(url string) bool {
	return videoPageRe.MatchString(url)
}

func isListPage(url string) bool {
	return listPageRe.MatchString(url)
}

func getVideosFromVideoPage(url string) (Course, error) {
	playListId, movieId := getVideoId(url)
	movie, err := getMoviesForAndroid(playListId)

	if err != nil {
		return nil, err
	}

	return &MovieResult{
		&Result{
			playListId: playListId,
			movie:      movie,
		},
		movieId,
	}, nil
}

func getVideosFromListPage(url string) (Course, error) {
	playListId, err := getVideoIdFromListPage(url)
	if err != nil {
		return nil, err
	}

	movie, err := getMoviesForAndroid(playListId)
	if err != nil {
		return nil, err
	}

	return &ListResult{
		&Result{
			playListId: playListId,
			movie:      movie,
		},
	}, nil
}

func getVideoId(url string) (string, string) {
	ids := videoIdRe.FindStringSubmatch(url)
	return ids[1], ids[2]
}

func getVideoIdFromListPage(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	match := playListIdRe.FindSubmatch(body)

	if len(match) < 2 {
		return "", errors.New("Can't find data-plid in " + url)
	}

	return string(match[1]), nil
}

func getMoviesForAndroid(playListId string) (*MovieData, error) {
	url := "http://mobile.open.163.com/movie/" + playListId + "/getMoviesForAndroid.htm"

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data MovieData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
