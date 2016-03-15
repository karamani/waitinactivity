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
			EnvVar:      "SEQREQUEST_DEBUG",
			Destination: &debugMode,
		},
		cli.IntFlag{
			Name:        "timeout",
			Usage:       "timeout in seconds",
			Destination: &timeoutArg,
		},
	}
	app.Action = func(c *cli.Context) {

		timeout := time.Second * time.Duration(timeoutArg)
		timer := time.NewTimer(timeout)

		streamActivity := make(chan bool)

		process := func(row []byte) error {
			debug(string(row))
			streamActivity <- true
			return nil
		}

		go func() {
			if err := iostreams.ReadStdin(process); err != nil {
				log.Panicln(err.Error())
			}
			os.Exit(1)
		}()

		for {
			select {
			case <-streamActivity:
				log.Println("activity")
			case <-timer.C:
				log.Println("timeout")
				fmt.Println("timeout")
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
