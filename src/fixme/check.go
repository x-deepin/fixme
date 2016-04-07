package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

var CMDCheck = cli.Command{
	Name:        "check",
	Usage:       "pid1 [pid2 ...]",
	Description: "Check whether the problems effected current system.",
	Action: func(c *cli.Context) {
		ids := c.Args()
		if len(ids) == 0 {
			cli.ShowCommandHelp(c, "check")
			return
		}
		force := c.Bool("force")

		ps, err := LoadProblems(c.GlobalString("db"))
		if err != nil || len(ps) == 0 {
			fmt.Println("E: The cache is empty. You need to run 'fixme update' first", err)
			return
		}
		for _, id := range ids {
			p := ps.Find(id)
			if p == nil {
				fmt.Println("Not found", id)
				continue
			}
			if force {
				p.Run(os.Stdout, "-c", "--force")
			} else {
				fmt.Println("Running...")
				fmt.Println("\n```")
				p.Run(os.Stdout, "-c", "-v")
				fmt.Println("```\n")
			}
		}
		if !force {
			fmt.Println("This project is developing, fix is default in dry-run mode.")
			fmt.Println("You can use -f to destroy your system :)")
		}
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force,f",
			Usage: "Do what i want.",
		},
	},
}
