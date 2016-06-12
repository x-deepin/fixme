package main

import (
	"archive/zip"
	"bufio"
	"crypto/md5"
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

var CMDUpdate = cli.Command{
	Name:        "update",
	Description: "Update remote fix scripts from $source to $cache",
	Action:      updateAction,
	Flags:       []cli.Flag{},
}

func updateAction(c *cli.Context) error {
	dest := c.GlobalString("cache")
	baseUrl := c.GlobalString("source")

	fingerprint, err := remoteFingerprint(baseUrl)
	if err != nil {
		return err
	}

	os.MkdirAll(dest, 0755)

	err = updateCache(baseUrl, fingerprint, dest)
	if err != nil {
		return err
	}

	fmt.Printf("Updated to newest PSet.\nYou can use \"fixme show\" to check the result.\n")
	return nil
}

func updateCache(baseUrl string, fingerprint string, destDir string) error {
	const NAME = "master.zip"
	fpath := path.Join(destDir, NAME)
	f, err := os.Open(fpath)
	if err == nil && checkFingerprint(f, fingerprint) {
		fmt.Printf("cache is newest --> %q\n", fingerprint)
		f.Close()
		return nil
	}

	f, err = os.Create(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, err := http.Get(baseUrl + "/" + fingerprint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tee := io.TeeReader(resp.Body, f)

	if !checkFingerprint(tee, fingerprint) {
		return fmt.Errorf("malform data")
	}

	err = uncompress(destDir, fpath)
	if err != nil {
		return err
	}
	return nil
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

func remoteFingerprint(baseUrl string) (string, error) {
	resp, err := http.Get(baseUrl + "/index")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	r := bufio.NewReader(resp.Body)

	_line, isPrefix, err := r.ReadLine()
	line := string(_line)
	if isPrefix {
		return line, fmt.Errorf("the fingerprint %q is too long", line)
	}
	return line, err
}

func checkFingerprint(r io.Reader, fingerprint string) bool {
	h := md5.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return false
	}
	if fingerprint != fmt.Sprintf("%x", h.Sum(nil)) {
		return false
	}
	return true
}
