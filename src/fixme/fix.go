package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

var CMDFix = cli.Command{
	Name:        "fix",
	Usage:       "pid1 [pid2 ...]",
	Description: "Try fixing the problems specified by pids",
	Action: func(c *cli.Context) {
		args := c.Args()
		if len(args) == 0 {
			cli.ShowCommandHelp(c, "fix")
			return
		}
		fmt.Printf("Fix %v .... Do nothing. (The function has't implement)\n", args)
	},
	Flags: []cli.Flag{},
}
