package main

import (
	"github.com/codegangsta/cli"
	"path"
)

var CMDCheck = cli.Command{
	Name:        "check",
	Usage:       "[pid1 pid2 ...]",
	Description: "Check whether the problems effected current system.",
	Action: func(c *cli.Context) error {
		db, err := LoadProblemDB(path.Join(c.GlobalString("cache"), ScriptDirName), c.GlobalString("db"))
		if err != nil {
			return err
		}

		var ps ProblemSet

		ids := c.Args()

		if len(ids) == 0 {
			ps = db.Search(func(Problem) bool { return true })
		} else {
			ps = db.Search(BuildSearchByIdFn(c.Args()))
		}

		ps.Run(Check)

		for _, p := range ps {
			db.Update(p)
		}
		return db.Save()
	},
}
