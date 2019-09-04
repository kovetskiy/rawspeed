package main

import (
	"fmt"
	"log"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/gdamore/tcell"
	"github.com/reconquest/stats-go"
)

var (
	version = "[manual build]"
	usage   = "layout " + version + `

Usage:
  rawspeed [options]

Options:
  --version                Show version.
`
)

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		panic(err)
	}

	_ = args

	up, err := watchKeyPress()
	if err != nil {
		log.Fatalln(err)
	}

	rate1s := stats.NewRate(time.Second)
	rate10s := stats.NewRate(time.Second * 10)

	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	err = screen.Init()
	if err != nil {
		panic(err)
	}

	quit := make(chan struct{})
	go func() {
		for {
			event := screen.PollEvent()
			switch event := event.(type) {
			case *tcell.EventKey:
				if event.Key() == tcell.KeyCtrlC {
					close(quit)
				}
			}
		}
	}()

	//screen.Show()

	width := getTerminalWidth()

	max1s := 0
	max10s := 0.0

loop:
	for {
		select {
		case <-up:
			rate1s.Increase()
			rate10s.Increase()

			speed1s := rate1s.Get()
			speed10s := float64(rate10s.Get()) / 10

			if speed1s > max1s {
				max1s = speed1s
			}
			if speed10s > max10s {
				max10s = speed10s
			}

			line := fmt.Sprintf(
				"k/s (1s): %v (avg 10s): %.2f | max (1s): %v (avg 10s): %.2f",
				speed1s,
				speed10s,
				max1s,
				max10s,
			)

			fmt.Print("\r" + line + getLineSuffix(width, len(line)))
		case <-quit:
			break loop
		}
	}

	screen.Fini()
}
