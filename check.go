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
		force := c.Bool("force")

		db, err := NewProblemDB(c.GlobalString("db"))
		if err != nil || len(db.cache) == 0 {
			return fmt.Errorf("The cache is empty. You need to run 'fixme update' first: %v", err)
		}
		for _, id := range ids {
			p := db.Find(id)
			if p == nil {
				fmt.Println("Not found", id)
				continue
			}
			if force {
				if !p.Check() {
					fmt.Printf("Found problem of %q\n", p.Id)
				}
				db.Add(p)
				db.Save()
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
		return nil
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force,f",
			Usage: "Do what i want.",
		},
	},
}
