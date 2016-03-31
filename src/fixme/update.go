package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

var CMDUpdate = cli.Command{
	Name:        "update",
	Usage:       "list all knowned problems",
	Description: "What is description?",
	Action:      updateAction,
	Flags:       []cli.Flag{},
}

func updateAction(c *cli.Context) {
	var ps ProblemSet
	print("\n")
	for i := 0; i < 10; i++ {
		ps = append(ps, GenRandomProblem(i+1))
		fmt.Printf("\rDownloading the problems of %3s/%d", ps[i].Id, 10)
	}
	print("\n")
	fmt.Println("Done! You can run 'fixme show' to check current problems.")
	SaveProblems(c.GlobalString("db"), ps)
}
