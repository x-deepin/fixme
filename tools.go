package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
)

func checkFileChanged(fpath string, fingerprint string) bool {
	f, err := os.Open(fpath)
	if err != nil {
		return true
	}
	defer f.Close()

	return !sameFingerprint(f, fingerprint)
}

func downloadFile(url string, dest string, fingerprint string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	tee := io.TeeReader(resp.Body, f)

	if !sameFingerprint(tee, fingerprint) {
		return fmt.Errorf("malform data")
	}

	return nil
}

func sameFingerprint(r io.Reader, fingerprint string) bool {
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

func remoteCatLine(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	r := bufio.NewReader(resp.Body)

	_line, isPrefix, err := r.ReadLine()
	line := string(_line)
	if isPrefix {
		return line, fmt.Errorf("the line %q is too long", line)
	}
	return line, err
}

func RED(s string) string {
	return "\033[31m" + s + "\033[0m"
}
