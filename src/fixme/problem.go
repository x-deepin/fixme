package main

import (
	"encoding/json"
	"fmt"
	"github.com/apcera/termtables"
	"os"
	"path"
)

const (
	EffectUnknown = "?"
	EffectYes     = "yes"
	EffectNo      = "no"
)
const (
	ScriptFix    = "fix"
	ScriptCheck  = "check"
	ScriptDetect = "detect"
)

type ProblemSet []Problem

type Problem struct {
	Id          string
	Title       string
	Description string

	Effected string

	dir string
}

func (ps ProblemSet) RenderSumary() string {
	t := termtables.CreateTable()
	t.AddHeaders("ID", "Title", "EffectMe")
	for _, p := range ps {
		t.AddRow(p.Id, p.Title, p.Effected)
	}
	return t.Render()
}

func (ps ProblemSet) Find(id string) *Problem {
	for _, p := range ps {
		if p.Id == id {
			return &p
		}
	}
	return nil
}

func (p Problem) String() string {
	return fmt.Sprintf("ID: %s\nTitle: %s\nDesc: %s\nEffectMe: %v\n",
		p.Id, p.Title, p.Description, p.Effected,
	)
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

func (p Problem) Fix() (string, error) {
	return ShellCode(path.Join(p.dir, ScriptCheck))
}
