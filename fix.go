package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

var CMDFix = cli.Command{
	Name:        "fix",
	Usage:       "pid1 [pid2 ...]",
	Description: "Try fixing the problems specified by pids",
	Action: func(c *cli.Context) error {
		ids := c.Args()
		if len(ids) == 0 {
			cli.ShowCommandHelp(c, "fix")
			return fmt.Errorf("Hasn't any pid")
		}
		dryRun := c.Bool("dry-run")

		db, err := LoadProblemDB(c.GlobalString("cache"))
		if err != nil {
			return err
		}
		for _, id := range ids {
			p := db.Find(id)
			if p == nil {
				fmt.Println("Not found", id)
				continue
			}

			if dryRun {
				fmt.Println("Running...")
				fmt.Println("\n```")
				p.Run(os.Stdout, "-f", "-v")
				fmt.Println("```\n")
			} else {
				p.Run(os.Stdout, "-f", "--force")
			}

		}
		return nil
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "dry-run,d",
			Usage: "Do what i want.",
		},
	},
}
