package fileidentifier

import (
	"math/big"
	"os/exec"
	"strings"
	"testing"
)

func getFileIDFromCommandExByName(t *testing.T, filename string) *big.Int {
	out, err := exec.Command("fsutil", "file", "queryfileid", filename).Output()
	if err != nil {
		t.Errorf("exec.Command(fsutil file queryfileid %v).Output() failed: %v", filename, err)
	}
	splitted := strings.Split(string(out), "0x")
	if len(splitted) != 2 {
		t.Errorf("out (%s) is no in the correct format, expected someting like: 'Datei-ID: 0x000000000000000000030000000618a1'", out)
	}
	expectedFileID := new(big.Int)
	parseSource := strings.TrimRight(splitted[1], "\r\n")
	_, ok := expectedFileID.SetString(parseSource, 16)
	if !ok || len(splitted) != 2 {
		t.Logf("filename: %q", filename)
		t.Logf("fsutil id output: %q", string(out))
		t.Logf("splitted[0]: %q", splitted[0])
		t.Logf("splitted[1]: %q", splitted[1])
		t.Logf("parseSource: %q", parseSource)
		t.Logf("fsutil id parsed: %s", expectedFileID.String())
		if !ok {
			t.FailNow()
		}
	}
	return expectedFileID
}
