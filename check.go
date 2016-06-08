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
	Action: func(c *cli.Context) error {
		ids := c.Args()
		if len(ids) == 0 {
			cli.ShowCommandHelp(c, "check")
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
				p.Run(os.Stdout, "-c", "-v")
				fmt.Println("```\n")

			} else {
				if !p.Check() {
					fmt.Printf("Found problem of %q\n", p.Id)
				}
				db.Add(p)
				db.Save()
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
