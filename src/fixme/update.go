package main

import (
	"archive/zip"
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

var CMDUpdate = cli.Command{
	Name:        "update",
	Usage:       "list all knowned problems",
	Description: "What is description?",
	Action:      updateAction,
	Flags:       []cli.Flag{},
}

func uncompressPSet(zipFile string, dest string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	write := func(name string, z *zip.File) error {
		os.MkdirAll(path.Dir(name), 0755)
		f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer f.Close()

		r, err := z.Open()
		if err != nil {
			return err
		}
		defer r.Close()

		_, err = io.Copy(f, r)
		return err
	}
	pathShift := func(name string) string {
		fs := strings.Split(name, "/")
		if len(fs) > 1 {
			return strings.Join(fs[1:], string(os.PathSeparator))
		}
		return strings.Join(fs, string(os.PathSeparator))
	}

	for _, f := range r.File {
		name := pathShift(f.Name)
		dname := path.Join(dest, name)

		if name == "functions" {
			err = write(dname, f)
			if err != nil {
				return err
			}
			continue
		}

		if f.FileInfo().IsDir() || path.Base(name) != ScriptFix {
			continue
		}

		err := write(dname, f)
		if err != nil {
			return err
		}
	}
	return nil
}

func ParsePSet(dest string) ([]*Problem, error) {
	var getId func(string) []string

	getId = func(dir string) []string {
		var r []string
		fs, err := ioutil.ReadDir(dir)
		if err != nil {
			return r
		}
		for _, f := range fs {
			if f.IsDir() {
				r = append(r, getId(path.Join(dir, f.Name()))...)
			} else if f.Name() == ScriptFix {
				r = append(r, path.Join(dir, f.Name()))
			}
		}
		return r
	}

	var ps []*Problem
	for _, id := range getId(dest) {
		p, err := NewProblem(dest, id)
		if err != nil {
			fmt.Println("E:", err)
			continue
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func SaveTo(url string, writer io.Writer) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(writer, resp.Body)
	return err
}

func downloadPSet(url string) (string, error) {
	f, err := ioutil.TempFile(os.TempDir(), "pset")
	if err != nil {
		return "", err
	}
	defer f.Close()

	err = SaveTo(url, f)
	if err != nil {
		return f.Name(), err
	}
	return f.Name(), nil
}

func updateAction(c *cli.Context) {
	pset, err := downloadPSet(c.GlobalString("pset"))
	if err != nil {
		fmt.Println("E:", err)
		return
	}
	defer os.Remove(pset)

	dest := c.GlobalString("cache")
	err = uncompressPSet(pset, dest)
	if err != nil {
		fmt.Println("E:", err)
		return
	}

	ps, err := ParsePSet(dest)
	if err != nil {
		fmt.Println("E:", err)
		return
	}

	db := &ProblemDB{
		dbPath: c.GlobalString("db"),
		cache:  make(map[string]*Problem),
	}
	for _, p := range ps {
		if p.AutoCheck {
			p.Check()
		}
		db.Add(p)
	}
	err = db.Save()
	if err != nil {
		fmt.Println("E:", err)
	}
	fmt.Println(db.RenderSumary())
}
