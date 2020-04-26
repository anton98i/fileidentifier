// +build !windows

package fileidentifier

import (
	"math/big"
	"os/exec"
	"strings"
	"testing"
)

func getFileIDFromCommandExByName(t *testing.T, path string) *big.Int {
	/* cool command, but osx can't handlt it
	  out, err := exec.Command("stat", "-c%i", path).Output()
		if err != nil {
			t.Errorf("exec.Command(stat -c%%i %v).Output() failed: %v", path, err)
	  }
	*/
	// -d shows the info of the directory instead of the file inside
	// -i shows the inode
	out, err := exec.Command("ls", "-d", "-i", path).Output()
	if err != nil {
		t.Errorf("exec.Command(ls -i %v).Output() failed: %v", path, err)
	}
	expectedFileID := new(big.Int)
	/*
		parseSource := strings.TrimRight(string(out), "\r\n")
	*/
	trimmed := strings.Trim(string(out), " \r\n")
	splitted := strings.Split(trimmed, " ")
	if len(splitted) != 2 {
		t.Logf("len(splitted) != 2, got: %d", len(splitted))
	}
	parseSource := splitted[0]
	_, ok := expectedFileID.SetString(parseSource, 10)
	if !ok || len(splitted) != 2 {
		t.Logf("filename: %q", path)
		t.Logf("stat id output: %q", string(out))
		t.Logf("trimmed: %q", trimmed)
		t.Logf("parseSource: %q", parseSource)
		t.Logf("stat id parsed: %q", expectedFileID.String())
		t.Logf("expectedFileID.SetString(%s, 16) failed", string(out))
		if !ok {
			t.FailNow()
		}
	}
	return expectedFileID
}
