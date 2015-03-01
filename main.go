package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/codegangsta/cli"
	"github.com/poying/necourse/necourse"
	"github.com/tj/go-spin"
)

var progName = filepath.Base(os.Args[0])

func main() {
	app := cli.NewApp()
	app.Name = progName
	app.Version = "0.0.0"
	cli.AppHelpTemplate = appHelpTemplate
	app.Action = action

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "concurrent,c",
			Value: 5,
			Usage: "同時下載數",
		},
		cli.StringFlag{
			Name:  "out,o",
			Value: ".",
			Usage: "檔案存放位置",
		},
	}

	app.Run(os.Args)
}

func action(c *cli.Context) {
	args := c.Args()

	if len(args) < 1 {
		cli.ShowAppHelp(c)
		return
	}

	//defer func() {
	//if err := recover(); err != nil {
	//fmt.Println("Error: ", err)
	//}
	//}()

	outputDir, err := filepath.Abs(c.String("out"))
	check(err)

	course, err := necourse.Get(args[0])
	check(err)

	downloader := NewDownloader(uint(c.Int("concurrent")))

	task := downloader.Task(course)

	downloader.Download(task, &Options{
		Quality:   necourse.SD,
		OutputDir: outputDir,
	})

	progress(course, task)
	task.Wait()
}

func progress(course necourse.Course, task *Task) {
	ticker := time.NewTicker(time.Millisecond * 100)
	spinner := spin.New()
	go func() {
		for range ticker.C {
			tick(course, task, spinner)
		}
	}()
}

func tick(course necourse.Course, task *Task, spinner *spin.Spinner) {
	spinFram := spinner.Next()
	iter := task.Status.Iter()

	fmt.Print("\033[0J\n")
	fmt.Printf("   \033[36m%s\033[m\n\n", course.Title())
	lines := 3

	for iter.HasNext() {
		video, status := iter.Next()
		lines += 1

		if status.Done() {
			if status.Failed() {
				fmt.Printf(" \033[31m✗\033[m %s (%s)\n", video.Title(), status.Error())
			} else {
				fmt.Printf(" \033[32m✓\033[m %s\n", video.Title())
			}
			continue
		}

		if status.Started() {
			fmt.Printf(" \033[36m%s\033[m %s \033[90m(%d bytes)\033[m\n", spinFram, video.Title(), status.Progress())
			continue
		}

		fmt.Printf("   \033[90m%s\033[m\n", video.Title())
	}

	fmt.Printf("\033[%dA", lines)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
