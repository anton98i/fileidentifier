package fileidentifier

import (
	"math/big"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func getFileIDFromCommandEx(t *testing.T, file *os.File) *big.Int {
	out, err := exec.Command("fsutil", "file", "queryfileid", file.Name()).Output()
	if err != nil {
		t.Errorf("exec.Command(fsutil file queryfileid %v).Output() failed: %v", file.Name(), err)
	}
	splitted := strings.Split(string(out), "0x")
	if len(splitted) != 2 {
		t.Errorf("out (%s) is no in the correct format, expected someting like: 'Datei-ID: 0x000000000000000000030000000618a1'", out)
		t.FailNow()
	}
	expectedFileID := new(big.Int)
	parseSource := strings.TrimRight(splitted[1], "\r\n")
	_, ok := expectedFileID.SetString(parseSource, 16)
	if !ok {
		t.Logf("filename: %q", file.Name())
		t.Logf("fsutil id output: %q", string(out))
		t.Logf("splitted[0]: %q", splitted[0])
		t.Logf("splitted[1]: %q", splitted[1])
		t.Logf("parseSource: %q", parseSource)
		t.Errorf("fsutil id parsed: %s", expectedFileID.String())
	}
	return expectedFileID
}
