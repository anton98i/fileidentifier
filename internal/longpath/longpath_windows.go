package longpath

import (
	"path/filepath"
	"strings"
)

const prefix = `\\?\`

// isPathShortEnough check if the path is short enough without prefix
func isPathShortEnough(path string) bool {
	return len(path) < 260
}

// Fix To avoid the default 260 character file path limit at windows is a prefix needed:\\?\
// e.g. \\?\c:\windows\foo.txt or \\?\UNC\server\foo\bar.txt.
// The extended form disables evaluation of . and .. path elements and disables the interpretation of / as equivalent to \.
// See https://msdn.microsoft.com/en-us/library/windows/desktop/aa365247(v=vs.85).aspx
//
// a error is only get returned, if the passed path is relative and is not possible resolveable to an absolute one
func Fix(path string) (string, error) {
	if isPathShortEnough(path) {
		return path, nil
	}

	if !filepath.IsAbs(path) {
		path = filepath.Clean(path)
		if isPathShortEnough(path) {
			// path got short enough after cleaning => return that
			return path, nil
		}

		abspath, err := filepath.Abs(path)
		if err != nil {
			return path, err
		}
		path = abspath
	}

	path = filepath.Clean(path)
	if isPathShortEnough(path) {
		// path got short enough after cleaning => return that
		return path, nil
	}

	if !strings.HasPrefix(path, prefix) {
		if strings.HasPrefix(path, `\\`) {
			// UNC paths e.g. \\server-name\foo\bar => \\?\UNC\server-name\foo\bar
			// remove first \ of path and add \\?\UNC
			path = prefix + "UNC" + path[1:]
		} else {
			path = prefix + path
		}
	}

	return path, nil
}
