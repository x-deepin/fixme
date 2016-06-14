package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"path"
)

var CMDList = cli.Command{
	Name:        "show",
	Usage:       "[pid ...]",
	Description: "List all known problems.",
	Action:      ActionShow,
	Flags:       []cli.Flag{},
}

func ActionShow(c *cli.Context) error {
	db, err := LoadProblemDB(path.Join(c.GlobalString("cache"), ScriptDirName), c.GlobalString("db"))
	if err != nil {
		return err
	}

	ids := c.Args()

	if len(ids) == 0 {
		fmt.Println(RED("Use \"fixme show\" $ID to check the detail information"))

		fmt.Println(db.RenderSumary())

		ps := db.Search(func(p Problem) bool { return p.Effected == EffectYes })
		if len(ps) != 0 {
			cmd := fmt.Sprintf("\tfixme fix")
			for _, p := range ps {
				cmd += fmt.Sprintf(" %s", p.Id)
			}

			fmt.Printf("\nTry using the command show below to fix the problmes\n\t%s\n",
				RED(cmd))
		}

		return nil
	} else {
		for _, p := range db.Search(BuildSearchByIdFn(c.Args())) {
			fmt.Println(p)
		}
	}

	return nil
}
