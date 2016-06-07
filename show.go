package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

var CMDList = cli.Command{
	Name:        "show",
	Usage:       "[pid ...]",
	Description: "List all known problems.",
	Action:      ActionShow,
	Flags:       []cli.Flag{},
}

func ActionShow(c *cli.Context) error {
	db, err := LoadProblemDB(c.GlobalString("cache"))
	if err != nil {
		return err
	}

	ids := c.Args()
	if len(ids) == 0 {
		fmt.Println(db.RenderSumary())
		return nil
	}

	for _, id := range ids {
		p := db.Find(id)
		if p == nil {
			fmt.Println("Not found", id)
			continue
		}

		fmt.Println(p)
	}
	return nil
}
