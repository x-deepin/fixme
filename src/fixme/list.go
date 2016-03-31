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

func ActionShow(c *cli.Context) {
	dbPath := c.GlobalString("db")
	ps, err := LoadProblems(dbPath)
	if err != nil || len(ps) == 0 {
		fmt.Println("E: The cache is empty. You need to run 'fixme update' first")
		return
	}

	ids := c.Args()
	if len(ids) == 0 {
		fmt.Println(ps.RenderSumary())
		return
	}

	for _, id := range ids {
		fmt.Println(ps.Render(id))
	}
}
