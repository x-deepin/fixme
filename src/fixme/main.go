package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()

	app.Version = "0.0.1"
	app.Name = "fixme"
	app.Usage = "Fix urgent bugs in deepin and eventually fix itself."
	app.Commands = []cli.Command{CMDList, CMDFix, CMDUpdate}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "server,s",
			Value: "http://fixme.deepin.com",
			Usage: "server url for updating and reporting",
		},
		cli.StringFlag{
			Name:  "db,d",
			Value: "db.json",
			Usage: "database path",
		},
	}

	app.Run(os.Args)
}
