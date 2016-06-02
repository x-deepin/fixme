package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()

	app.Version = "0.0.1"
	app.Name = "fixme"
	app.Usage = "Fix urgent bugs in deepin and eventually fix itself."
	app.Commands = []cli.Command{CMDList, CMDFix, CMDCheck, CMDUpdate}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "pset,s",
			Value: "https://github.com/x-deepin/p/archive/master.zip",
			Usage: "server url for updating and reporting",
		},
		cli.StringFlag{
			Name:  "cache,c",
			Value: "/dev/shm/p",
			Usage: "the cache directory to store pset scripts",
		},
		cli.StringFlag{
			Name:  "db,d",
			Value: "db.json",
			Usage: "database path",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "E:", err)
		os.Exit(-1)
	}
}
