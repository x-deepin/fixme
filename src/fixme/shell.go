package main

import (
	"io/ioutil"
	"os"
	"os/exec"
)

// ShellCode run the code by $SHELL with args and
// the stderr and stdout was combined.
func ShellCode(code string, args ...string) (string, error) {
	f, err := ioutil.TempFile(os.TempDir(), "shell_code")
	if err != nil {
		return "", err
	}
	defer f.Close()
	defer os.Remove(f.Name())

	ioutil.WriteFile(f.Name(), ([]byte)(code), 0755)

	args = append([]string{f.Name()}, args...)
	cmd := exec.Command(os.ExpandEnv("$SHELL"), args...)

	d, err := cmd.CombinedOutput()
	return string(d), err
}
