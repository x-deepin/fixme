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

type Problem struct {
	Id          string
	Title       string
	Description string
	ScriptPath  string

	Effected  string
	LastCheck time.Time

	AutoCheck bool   `json:"AUTO_CHECK"`
	AutoFix   bool   ` json:"AUTO_FIX"`
	Author    string `json:"AUTHOR"`
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
		Effected:   EffectUnknown,
	}
	var err error
	buf := bytes.NewBuffer(nil)

	err = p.Run(buf, "--verify")
	if buf.String() != "verified fixme script" {
		return nil, fmt.Errorf("Invalid script(%s):%v", fixPath, err)
	}
	buf.Reset()

	err = p.Run(buf, "-t")
	if err != nil {
		return nil, err
	}
	p.Title = strings.TrimSpace(buf.String())
	buf.Reset()

	p.Run(buf, "-m")
	err = json.Unmarshal(buf.Bytes(), &p)
	buf.Reset()
	return p, err
}

func (p *Problem) Check() bool {
	err := p.Run(os.Stdout, "-c", "--force")
	if err != nil {
		p.Effected = EffectYes
		return false
	} else {
		p.Effected = EffectNo
		return true
	}
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

type ProblemDB struct {
	dbPath string
	cache  map[string]*Problem
}

func (db ProblemDB) Add(p *Problem) {
	db.cache[p.Id] = p
}

func (db ProblemDB) Find(id string) *Problem {
	return db.cache[id]
}

func (db ProblemDB) RenderSummary() string {
	t := termtables.CreateTable()
	t.AddHeaders("ID", "Title", "EffectMe", "AutoCheck")
	for _, p := range db.cache {
		t.AddRow(p.Id, p.Title, p.Effected, p.AutoCheck)
	}
	return t.Render()
}

func (db ProblemDB) Save() error {
	f, err := os.Create(db.dbPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(db.cache)
}

func NewProblemDB(dbPath string) (*ProblemDB, error) {
	f, err := os.Open(dbPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	db := &ProblemDB{
		dbPath: dbPath,
		cache:  make(map[string]*Problem),
	}
	return db, json.NewDecoder(f).Decode(&db.cache)
}

func (db ProblemDB) build(sourceDir string) {
}
