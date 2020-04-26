package longpath

import (
	"path/filepath"
	"strings"
	"testing"
)

var veryLong string

func init() {
	veryLong = "l" + strings.Repeat("o", 248) + "ng"
}

func TestFixLongPath(t *testing.T) {
	for _, test := range []struct{ in, want string }{
		// Short; unchanged:
		{`C:\short.txt`, `C:\short.txt`},
		{`C:\`, `C:\`},
		{`C:`, `C:`},
		// The "long" substring is replaced by a looooooong with 248 o's
		{`C:\long\foo.txt`, `\\?\C:\long\foo.txt`},
		// should contert / to \
		{`C:/long/foo.txt`, `\\?\C:\long\foo.txt`},
		// should resolve .
		{`C:\long\foo\\bar\.\baz\\`, `\\?\C:\long\foo\bar\baz`},
		// \\ path
		{`\\unc\path`, `\\unc\path`},
		{`C:long.txt`, `C:long.txt`},
		// should resolve ..
		{`c:\long\..\bar\baz`, `c:\bar\baz`},
		// the result of this test is depending of the short limit:
		// {`C:\long\foo\\bar\..\..\baz\\`, `\\?\C:\long\baz`},
		{`C:\long\foo\\bar\..\..\baz\\`, `C:\long\baz`},
		// should resolve .. and is not short enough to return it without \\?\
		{`C:\long\foo\\bar\..\..\..\baz\\`, `C:\baz`},
		// already a extended path (return it unchanged)
		{`\\?\c:\long\foo.txt`, `\\?\c:\long\foo.txt`},
		{`\\?\c:\long/foo.txt`, `\\?\c:\long\foo.txt`},
		{`\\?\C:\long/foo.txt`, `\\?\C:\long\foo.txt`},
		{`\\?\D:\long/foo.txt`, `\\?\D:\long\foo.txt`},
		// unc paths
		{`\\`, `\\`},
		{`\\foo`, `\\foo`},
		{`\\foo\long\bar`, `\\?\UNC\foo\long\bar`},
		// should make relative paths absolute
		{`long.txt`, `long.txt`},
		{`long1111111111111111111111111111111111111111.txt`, `long1111111111111111111111111111111111111111.txt`},
		// a too long path gets short after the .. clean
		{`foo\long\long\..\..\bar.txt`, `foo\bar.txt`},
	} {
		testBasic(t, test.in, test.want)
		/* breakpoint in if for failed tets for easy debugging:
		if !(testBasic(t, test.in, test.want)) {
			testBasic(t, test.in, test.want)
		}
		*/
	}
}

func testBasic(t *testing.T, inSort, wantShort string) bool {
	success := true
	in := strings.ReplaceAll(inSort, "long", veryLong)
	want := strings.ReplaceAll(wantShort, "long", veryLong)
	if !isPathShortEnough(want) && !filepath.IsAbs(want) {
		want = filepath.Clean(want)
		if !isPathShortEnough(want) {
			var err error
			want, err = filepath.Abs(want)
			if err != nil {
				t.Errorf("filepath.Abs(%s) failed: %v", want, err)
				success = false
			}
			want = filepath.Clean(want)
			if !isPathShortEnough(want) {
				want = `\\?\` + want
			}
		}
	}
	_, expectedError := filepath.Abs(in)
	got, err := Fix(in)
	if err != expectedError {
		t.Errorf("Fix(%s) error: got: %v; expected: %v", inSort, err, expectedError)
		success = false
		if err != nil {
			// if a error happens => we want to receive the input cleaned
			want = filepath.Clean(in)
		}
	}
	if got != want {
		gotShortend := strings.ReplaceAll(got, veryLong, "long")
		wantShort := strings.ReplaceAll(want, veryLong, "long")
		t.Errorf("Fix(%s) got: %s; expected %s", inSort, gotShortend, wantShort)
		success = false
	}
	return success
}
