package main

import (
	"archive/zip"
	"bytes"
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

func tryReadContentInZip(f *zip.File) (string, error) {
	r, err := f.Open()
	if err != nil {
		return "", err
	}
	defer r.Close()
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, r)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func ParsePSet(zipFile string) (ProblemSet, error) {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return nil, err
	}

	cache := make(map[string]*zip.File)
	ids := make(map[string]struct{})

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if path.Base(f.Name) == ScriptFix {
			ids[path.Dir(f.Name)] = struct{}{}
		}
		cache[f.Name] = f
	}
	var ps ProblemSet
	for id := range ids {
		p := Problem{}
		p.Id = strings.Join(strings.Split(id, "/")[1:], ".")

		if f, ok := cache[path.Join(id, ScriptFix)]; ok {
			p.ContentFix, err = tryReadContentInZip(f)
			if err != nil {
				fmt.Println("W:", err)
			}
		}

		if f, ok := cache[path.Join(id, ScriptCheck)]; ok {
			p.ContentCheck, err = tryReadContentInZip(f)
			if err != nil {
				fmt.Println("W:", err)
			}
		}

		if f, ok := cache[path.Join(id, ScriptDetect)]; ok {
			p.ContentDetect, err = tryReadContentInZip(f)
			if err != nil {
				fmt.Println("W:", err)
			}
		}

		if f, ok := cache[path.Join(id, ScriptDescription)]; ok {
			p.Description, err = tryReadContentInZip(f)
			if err != nil {
				fmt.Println("W:", err)
			}
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

	ps, err := ParsePSet(pset)
	if err != nil {
		fmt.Println("E:", err)
		return
	}

	fmt.Println("Downloaded:", ps.RenderSumary())

	err = SaveProblems(c.GlobalString("db"), ps)
	if err != nil {
		fmt.Println("E:", err)
	}
}
