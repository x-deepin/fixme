package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func RED(s string) string {
	return "\033[31m" + s + "\033[0m"
}

func DoCheck(ps []*Problem, dryRun bool) error {
	for _, p := range ps {
		fmt.Printf("Checking problem \"%s\"\n", RED(p.Title))
		if dryRun {
			fmt.Println("\n```")
			p.Run(os.Stdout, "-c", "-v")
			fmt.Printf("```\n\n")
		} else {
			if !p.Check() {
				fmt.Printf("Found problem of %q\n", p.Id)
			}
		}
	}
	return nil
}

var CMDCheck = cli.Command{
	Name:        "check",
	Usage:       "[pid1 pid2 ...]",
	Description: "Check whether the problems effected current system.",
	Action: func(c *cli.Context) error {
		db, err := LoadProblemDB(c.GlobalString("cache"), c.GlobalString("db"))
		if err != nil {
			return err
		}

		var ps []*Problem

		ids := c.Args()

		if len(ids) == 0 {
			ps = db.List()
		} else {
			for _, id := range ids {
				p := db.Find(id)
				if p == nil {
					fmt.Println("Not found", id)
					continue
				}
				ps = append(ps, p)
			}
		}

		DoCheck(ps, c.Bool("dry-run"))
		for _, p := range ps {
			db.Update(p)
		}
		return db.Save()
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "dry-run,d",
			Usage: "Do what i want.",
		},
	},
}
