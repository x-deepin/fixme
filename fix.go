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
			err := p.Fix()
			if err != nil {
				return err
			}
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

		ids := c.Args()

		if len(ids) == 0 {
			ps = db.Search(func(p Problem) bool {
				return p.Effected == EffectYes || p.Effected == EffectUnknown
			})
		} else {
			ps = db.Search(BuildSearchByIdFn(c.Args()))
		}

		err = DoFix(ps, c.Bool("dry-run"))
		if err != nil {
			return err
		}

		return db.Save()
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "dry-run,d",
			Usage: "Do what I want.",
		},
	},
}
