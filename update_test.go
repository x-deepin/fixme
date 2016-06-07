package main

import (
	"testing"
)

func TestRemoteFingerprint(t *testing.T) {
	line, err := remoteFingerprint("http://packages.deepin.com/deepin/fixme")
	if err != nil {
		t.Fatal("E:", err)
	}
	if len(line) != 32 {
		t.Fatalf("%q not a md5sum value", line)
	}
}
