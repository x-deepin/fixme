package main

import (
	"github.com/codegangsta/cli"
	"path"
)

var CMDFix = cli.Command{
	Name:        "fix",
	Usage:       "pid1 [pid2 ...]",
	Description: "Try fixing the problems specified by pids",
	Action: func(c *cli.Context) error {
		db, err := LoadProblemDB(path.Join(c.GlobalString("cache"), ScriptDirName), c.GlobalString("db"))
		if err != nil {
			return err
		}

		var ps ProblemSet

		ids := c.Args()

		if len(ids) == 0 {
			ps = db.Search(func(p Problem) bool {
				return p.Effected == EffectYes || p.Effected == EffectUnknown
			})
		} else {
			ps = db.Search(BuildSearchByIdFn(c.Args()))
		}

		err = ps.Run(Fix)
		if err != nil {
			return err
		}

		return db.Save()
	},
}
