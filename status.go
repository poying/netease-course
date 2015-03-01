package main

import (
	"io"

	"github.com/poying/necourse/necourse"
)

type VideoStatus struct {
	reader   io.Reader
	video    *necourse.Video
	bytes    int
	started  bool
	finished bool
	err      error
}

func (s *VideoStatus) Progress() int {
	return s.bytes
}

func (s *VideoStatus) Started() bool {
	return s.started
}

func (s *VideoStatus) Done() bool {
	return s.finished
}

func (s *VideoStatus) Failed() bool {
	return s.err != nil
}

func (s *VideoStatus) Error() error {
	return s.err
}

func (s *VideoStatus) Read(p []byte) (int, error) {
	n, err := s.reader.Read(p)
	s.bytes += n
	return n, err
}

type Status struct {
	videos      []necourse.Video
	videoStatus map[string]*VideoStatus
}

func NewStatus(videos []necourse.Video) *Status {
	videoStatus := make(map[string]*VideoStatus)
	size := len(videos)

	for index := 0; index < size; index += 1 {
		video := &videos[index]
		videoStatus[video.Id()] = &VideoStatus{
			video:    video,
			bytes:    0,
			finished: false,
		}
	}

	return &Status{
		videos,
		videoStatus,
	}
}

func (s *Status) VideoStatus(video *necourse.Video) *VideoStatus {
	status, ok := s.videoStatus[video.Id()]
	if ok {
		return status
	}
	return nil
}

func (s *Status) Iter() *StatusIterator {
	return &StatusIterator{
		status: s,
		index:  0,
	}
}

type StatusIterator struct {
	status *Status
	index  int
}

func (si *StatusIterator) Next() (*necourse.Video, *VideoStatus) {
	index := si.index
	si.index += 1
	video := &si.status.videos[index]
	status := si.status.VideoStatus(video)
	return video, status
}

func (si *StatusIterator) HasNext() bool {
	return si.index < len(si.status.videos)
}
