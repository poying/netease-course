package main

import (
	"io"
	"net/http"
	"os"
	"fmt"
	"path"
	"strings"
	"strconv"
	"path/filepath"

	"github.com/poying/necourse/necourse"
)

var pathSeparator = strconv.QuoteRune(os.PathSeparator)

type Downloader struct {
	concurrent uint
}

func NewDownloader(concurrent uint) *Downloader {
	return &Downloader{concurrent}
}

func (d *Downloader) Download(task *Task, opts *Options) {
	in := task.Channel
	sem := make(chan bool, d.concurrent)

	go func() {
		defer close(sem)
		for video := range in {
			sem <- true
			go func(video *necourse.Video) {
				defer func() {
					<-sem
					task.Done()
				}()
				d.downloadVideo(video, task.Status.VideoStatus(video), opts)
			}(video)
		}
	}()
}

func (d *Downloader) Task(course necourse.Course) *Task {
	videos := d.getVideos(course)
	task := &Task{
		Channel: d.makeVideoChan(videos),
		Status:  NewStatus(videos),
		course:  course,
	}
	task.Add(len(videos))
	return task
}

func (d *Downloader) makeVideoChan(videos []necourse.Video) <-chan *necourse.Video {
	out := make(chan *necourse.Video, len(videos))
	size := len(videos)

	for index := 0; index < size; index += 1 {
		out <- &videos[index]
	}

	close(out)
	return out
}

func (d *Downloader) getVideos(course necourse.Course) []necourse.Video {
	if movie, ok := course.(*necourse.MovieResult); ok {
		video := movie.Video()
		if video == nil {
			return make([]necourse.Video, 0)
		}
		list := []necourse.Video{*video}
		return list
	}

	return course.Videos()
}

func (d *Downloader) downloadVideo(video *necourse.Video, status *VideoStatus, opts *Options) {
	defer func() { status.finished = true }()

	url := video.Url(opts.Quality)
	extname := filepath.Ext(url)
	outputFile := path.Join(opts.OutputDir, videoFileName(video)) + extname

	writer, err := os.Create(outputFile)
	if err != nil {
		status.err = err
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		status.err = err
		return
	}
	defer resp.Body.Close()
	status.reader = resp.Body
	status.started = true

	_, err = io.Copy(writer, status)
	if err != nil {
		status.err = err
	}
}

func videoFileName(video *necourse.Video) string {
	number := video.PNumber()
	title := strings.Replace(video.Title(), pathSeparator, "_", -1)
	return fmt.Sprintf("%d. %s", number, title)
}
