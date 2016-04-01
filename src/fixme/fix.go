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
		ids := c.Args()
		if len(ids) == 0 {
			cli.ShowCommandHelp(c, "fix")
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
			fmt.Println("Running...")
			fmt.Println("\n```")
			fmt.Println(p.ContentFix)
			fmt.Println("```\n")

			if force {
				fmt.Println(p.Fix())
			}
		}
		fmt.Println("This project is developing, fix is default in dry-run mode.")
		fmt.Println("You can use -f to destroy your system :)")
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force,f",
			Usage: "Do what i want.",
		},
	},
}
