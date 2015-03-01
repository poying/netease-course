package main

import (
	"sync"

	"github.com/poying/necourse/necourse"
)

type Task struct {
	sync.WaitGroup
	Channel <-chan *necourse.Video
	Status  *Status
	course  necourse.Course
}
