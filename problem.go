package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/apcera/termtables"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sort"
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

type ProblemSet []*Problem

func (ps ProblemSet) Less(i, j int) bool {
	if ps[i].Effected == ps[j].Effected {
		return ps[i].Id > ps[j].Id
	}
	return ps[i].Effected > ps[j].Effected
}
func (ps ProblemSet) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}
func (ps ProblemSet) Len() int {
	return len(ps)
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

	p.Title = strings.Trim(buf.String(), " \n\r")
	buf.Reset()

	p.Run(buf, "-m")
	err = json.Unmarshal(buf.Bytes(), &p)
	buf.Reset()
	return p, err
}

func (p *Problem) Check() bool {
	err := p.Run(os.Stdout, "-c", "--force")
	p.LastCheck = time.Now()
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

func (db ProblemDB) Update(p *Problem) error {
	if _, ok := db.cache[p.Id]; !ok {
		return fmt.Errorf("There hasn't %q problem in db", p.Id)
	}
	db.cache[p.Id] = p
	return nil
}

func (db ProblemDB) sort() []*Problem {
	var r ProblemSet

	for _, p := range db.cache {
		r = append(r, p)
	}
	sort.Sort(r)
	return r
}

func (db ProblemDB) RenderSumary() string {
	t := termtables.CreateTable()
	t.AddHeaders("ID", "Title", "LastCheck")
	var ps []string
	for _, p := range db.sort() {
		lc := p.LastCheck.Format("2006-01-02 15:04:05")
		if p.LastCheck.IsZero() {
			lc = "never"
		}
		if p.Effected == EffectYes {
			ps = append(ps, p.Id)
			t.AddRow(RED(p.Id), p.Title, lc)
		} else {
			t.AddRow(p.Id, p.Title, lc)
		}
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

func scanProblemIDs(dir string) []string {
	var r []string
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return r
	}
	for _, f := range fs {
		if f.IsDir() {
			r = append(r, scanProblemIDs(path.Join(dir, f.Name()))...)
		} else if f.Name() == ScriptFix {
			r = append(r, path.Join(dir, f.Name()))
		}
	}
	return r

}
func BuildProblemDB(cacheDir string, dbPath string) (*ProblemDB, error) {
	var ps []*Problem
	for _, id := range scanProblemIDs(cacheDir) {
		p, err := NewProblem(cacheDir, id)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}

	os.MkdirAll(path.Dir(dbPath), 0755)
	db := &ProblemDB{
		dbPath: dbPath,
		cache:  make(map[string]*Problem),
	}

	for _, p := range ps {
		if p.AutoCheck {
			p.Check()
		}
		db.cache[p.Id] = p
	}
	fmt.Printf(RED("BuildProblemDB from %q to %q\n"), cacheDir, dbPath)
	return db, db.Save()
}

func LoadProblemDB(cacheDir string, dbPath string) (*ProblemDB, error) {
	f, err := os.Open(dbPath)
	if err != nil {
		return BuildProblemDB(cacheDir, dbPath)
	}
	defer f.Close()

	db := &ProblemDB{
		dbPath: dbPath,
		cache:  make(map[string]*Problem),
	}

	err = json.NewDecoder(f).Decode(&db.cache)

	if len(db.cache) == 0 {
		db, err = BuildProblemDB(cacheDir, dbPath)
	}

	if err != nil || len(db.cache) == 0 {
		return nil, fmt.Errorf("The cache is empty. You need to run 'fixme update' first: %v", err)
	}
	fmt.Printf(RED("LoadProblemDB from %q\n"), dbPath)
	return db, nil
}

func BuildSearchByIdFn(ids []string) func(p Problem) bool {
	return func(p Problem) bool {
		for _, id := range ids {
			if p.Id == id {
				return true
			}
		}
		return false
	}
}

type SearchFn func(Problem) bool

func (db ProblemDB) Search(fn SearchFn) []*Problem {
	var r []*Problem
	for _, p := range db.sort() {
		if fn(*p) {
			r = append(r, p)
		}
	}
	return r
}
