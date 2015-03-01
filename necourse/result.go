package necourse

import "strings"

type Course interface {
	PlayListId() string
	Description() string
	Director() string
	ImgUrl() string
	LargeImgUrl() string
	School() string
	Source() string
	SubTitle() string
	Title() string
	Tags() []string
	Videos() []Video
	VideoCount() int
}

type Result struct {
	playListId string
	movie      *MovieData
}

func (r *Result) PlayListId() string {
	return r.playListId
}

func (r *Result) Description() string {
	return r.getMovie().Description
}

func (r *Result) Director() string {
	return r.getMovie().Director
}

func (r *Result) ImgUrl() string {
	return r.getMovie().ImgPath
}

func (r *Result) LargeImgUrl() string {
	return r.getMovie().LargeImgUrl
}

func (r *Result) School() string {
	return r.getMovie().School
}

func (r *Result) Source() string {
	return r.getMovie().Source
}

func (r *Result) Title() string {
	return r.getMovie().Title
}

func (r *Result) SubTitle() string {
	return r.getMovie().SubTitle
}

func (r *Result) Tags() []string {
	tags := r.getMovie().Tags
	return strings.Split(tags, ",")
}

func (r *Result) Videos() []Video {
	movie := r.getMovie()
	videos := make([]Video, len(movie.VideoList))
	for index, video := range movie.VideoList {
		videos[index] = Video{video}
	}
	return videos
}

func (r *Result) VideoCount() int {
	return len(r.movie.VideoList)
}

func (r *Result) getMovie() *MovieData {
	return r.movie
}

type MovieResult struct {
	*Result
	id string
}

func (m *MovieResult) Id() string {
	return m.id
}

func (m *MovieResult) Video() *Video {
	for _, video := range m.Videos() {
		if video.Id() == m.Id() {
			return &video
		}
	}
	return nil
}

type ListResult struct {
	*Result
}
