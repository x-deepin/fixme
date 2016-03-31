package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

var CMDFix = cli.Command{
	Name:        "fix",
	Usage:       "pid1 [pid2 ...]",
	Description: "Try fixing the problems specified by pid",
	Action: func(c *cli.Context) {
		fmt.Println("Fix %v .... Do nothing. (The function Hasn't implement)", c.Args())
	},
	Flags: []cli.Flag{},
}
