package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func DoFix(ps []*Problem, dryRun bool) error {
	for _, p := range ps {
		if dryRun {
			fmt.Println("Running...")
			fmt.Println("\n```")
			p.Run(os.Stdout, "-f", "-v")
			fmt.Printf("```\n\n")
		} else {
			p.Run(os.Stdout, "-f", "--force")
		}
	}
	return nil
}

var CMDFix = cli.Command{
	Name:        "fix",
	Usage:       "pid1 [pid2 ...]",
	Description: "Try fixing the problems specified by pids",
	Action: func(c *cli.Context) error {
		db, err := LoadProblemDB(c.GlobalString("cache"), c.GlobalString("db"))
		if err != nil {
			return err
		}

		var ps []*Problem

		if c.Bool("autofix") {
			ps = db.Search(func(p Problem) bool {
				return p.AutoFix
			})
		} else {
			ps = db.Search(BuildSearchByIdFn(c.Args()))
			if len(ps) == 0 {
				cli.ShowCommandHelp(c, "fix")
				return fmt.Errorf("Hasn't any pid")
			}
		}
		return DoFix(ps, c.Bool("dry-run"))
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "dry-run,d",
			Usage: "Do what I want.",
		},
		cli.BoolFlag{
			Name:  "autofix",
			Usage: "fix all script which AUTO_FIX==true",
		},
	},
}
