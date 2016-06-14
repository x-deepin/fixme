package main

import (
	"archive/zip"
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"os"
	"path"
	"strings"
)

const ScriptDirName = "scripts"

var CMDUpdate = cli.Command{
	Name:        "update",
	Description: "Update remote fix scripts from $source to $cache",
	Action:      updateAction,
	Flags:       []cli.Flag{},
}

func updateAction(c *cli.Context) error {
	cacheDir := c.GlobalString("cache")
	baseUrl := c.GlobalString("source")

	fingerprint, err := remoteCatLine(baseUrl + "/index")
	if err != nil {
		return err
	}

	os.MkdirAll(cacheDir, 0755)

	err = updateCache(baseUrl, fingerprint, cacheDir)
	if err != nil {
		return err
	}

	fmt.Printf("Updated to newest PSet.\nYou can use \"fixme show\" to check the result.\n")
	return nil
}

func updateCache(baseUrl string, fingerprint string, cacheDir string) error {
	const NAME = "master.zip"
	zipFile := path.Join(cacheDir, NAME)

	if !checkFileChanged(zipFile, fingerprint) {
		fmt.Printf("cache is newest --> %q\n", fingerprint)
		return nil
	}

	err := downloadFile(baseUrl+"/"+fingerprint, zipFile, fingerprint)
	if err != nil {
		return err
	}

	scriptDir := path.Join(cacheDir, ScriptDirName)
	os.RemoveAll(scriptDir)
	return uncompress(scriptDir, zipFile)
}

func uncompress(destDir string, zipFile string) error {
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
		dname := path.Join(destDir, name)

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
