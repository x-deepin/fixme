package main

import (
	"encoding/json"
	"fmt"
	"github.com/apcera/termtables"
	"os"
)

const (
	EffectUnknown = "?"
	EffectYes     = "yes"
	EffectNo      = "no"
)

type ProblemSet []Problem

type Problem struct {
	Id          string
	Title       string
	Description string

	Effected string

	fixScript    string
	checkScript  string
	detectScript string
}

func (ps ProblemSet) RenderSumary() string {
	t := termtables.CreateTable()
	t.AddHeaders("ID", "Title", "EffectMe")
	for _, p := range ps {
		t.AddRow(p.Id, p.Title, p.Effected)
	}
	return t.Render()
}

func (ps ProblemSet) Render(id string) string {
	for _, p := range ps {
		if p.Id != id {
			continue
		}
		return fmt.Sprintf("ID: %s\nTitle: %s\nDesc: %s\nEffectMe: %v\n",
			p.Id, p.Title, p.Description, p.Effected,
		)
	}
	return "not found " + id + "\n"
}

func SaveProblems(fpath string, ps ProblemSet) error {
	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(ps)
}

func LoadProblems(fpath string) (ProblemSet, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var r ProblemSet
	err = json.NewDecoder(f).Decode(&r)
	return r, err
}
