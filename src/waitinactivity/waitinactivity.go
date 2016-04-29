package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/karamani/iostreams"
)

var (
	debugMode  bool
	timeoutArg int
)

func main() {

	app := cli.NewApp()
	app.Name = "waitinactive"
	app.Usage = "Wait inactive"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "debug mode",
			Destination: &debugMode,
		},
		cli.IntFlag{
			Name:        "timeout",
			Usage:       "timeout in seconds",
			Value:       1,
			Destination: &timeoutArg,
		},
	}
	app.Action = func(c *cli.Context) {

		timeout := time.Second * time.Duration(timeoutArg)
		timer := time.NewTimer(timeout)

		chanEOF := make(chan bool)
		streamActivity := make(chan []byte)

		go func() {
			if err := iostreams.ChanStdin(streamActivity); err != nil {
				log.Panicln(err.Error())
			}
			chanEOF <- true
		}()

		for {
			select {
			case msg := <-streamActivity:
				debug(string(msg))
				log.Println("activity")
			case <-timer.C:
				log.Println("timeout")
				fmt.Println("timeout")
			case <-chanEOF:
				log.Println("eof")
				return
			}
			timer.Reset(timeout)
		}
	}

	app.Run(os.Args)
}

func debug(msg string) {
	if debugMode {
		log.Println(msg)
	}
}
