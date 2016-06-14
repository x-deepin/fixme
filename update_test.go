package main

import (
	"testing"
)

func TestRemoteFingerprint(t *testing.T) {
	line, err := remoteCatLine("http://packages.deepin.com/deepin/fixme/index")
	if err != nil {
		t.Fatal("E:", err)
	}
	if len(line) != 32 {
		t.Fatalf("%q not a md5sum value", line)
	}
}
