package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/apcera/termtables"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

const (
	EffectUnknown = "?"
	EffectYes     = "yes"
	EffectNo      = "no"
)
const (
	ScriptFix         = "fix"
	ScriptCheck       = "check"
	ScriptDetect      = "detect"
	ScriptDescription = "README.md"
)

type ProblemSet []*Problem

type Problem struct {
	Id          string
	Title       string
	Description string
	ScriptPath  string

	Effected  string
	LastCheck time.Time
	AutoCheck bool
}

func (ps ProblemSet) RenderSumaryTest() string {
	// TODO: parse the README.md contents
	var r string
	for _, p := range ps {
		r = r + p.String() + "\n\n"
	}
	return r
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
			return p
		}
	}
	return nil
}

func NewProblem(base, fixPath string) (*Problem, error) {
	if !strings.HasPrefix(fixPath, base) || base == fixPath {
		return nil, fmt.Errorf("Invalid fixPath:%v %v", base, fixPath)
	}

	id := strings.Replace(fixPath, base, "", 1)
	id = strings.Replace(id, string(os.PathSeparator), ".", -1)
	id = id[0 : len(id)-len("fix")]
	id = strings.Trim(id, ".")

	p := &Problem{
		Id:         id,
		ScriptPath: fixPath,
	}
	var err error
	buf := bytes.NewBuffer(nil)
	err = p.Run(buf, "-t")
	if err != nil {
		return nil, err
	}
	p.Title = strings.TrimSpace(buf.String())
	return p, nil
}

func (p Problem) Run(output io.Writer, arg ...string) error {
	cmd := exec.Command(p.ScriptPath, arg...)
	cmd.Dir = path.Dir(p.ScriptPath)
	cmd.Stdout = output
	err := cmd.Run()
	return err
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
