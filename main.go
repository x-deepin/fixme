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
	app.Usage = "Fix urgent bugs in deepin and eventually fix itself. See also https://github.com/x-deepin/p"

	app.Commands = []cli.Command{CMDList, CMDFix, CMDCheck, CMDUpdate}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "source,s",
			Value: "http://cdn.packages.deepin.com/deepin/fixme",
			Usage: "server url for updating and reporting",
		},
		cli.StringFlag{
			Name:  "cache,c",
			Value: "/var/lib/fixme",
			Usage: "the cache directory to store scripts",
		},
		cli.StringFlag{
			Name:  "db",
			Value: os.ExpandEnv("$HOME/.cache/fixme/db.json"),
			Usage: "the cache directory to store scripts",
		},
	}

	var err error
	if len(os.Args) == 1 {
		err = app.Run(append(os.Args, "show"))
	} else {
		err = app.Run(os.Args)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "E:", err)
		os.Exit(-1)
	}
}
