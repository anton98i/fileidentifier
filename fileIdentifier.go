package fileidentifier

import (
	"fmt"
	"math/big"
	"os"
)

// inspiered mostly by if not all from: https://git.icinga.com/github-mirror/icingabeat/-/tree/811e7afc940fd0c332242dcd8c20a89b63f3d19f/vendor/github.com/elastic/beats/filebeat/input/file
// just added big.Int

// getBigInt returns num << n
func getBigInt(num uint64, n uint) *big.Int {
	var n1 big.Int
	n1.SetUint64(num)
	// != 0 check is done inside the shl function
	// Lsh: https://golang.org/src/math/big/int.go?s=25314:25352#L993
	// shl: https://golang.org/src/math/big/nat.go#L981
	return n1.Lsh(&n1, n)
}

// getBigIntRsh returns num >> n
func getBigIntRsh(num uint64, n uint) *big.Int {
	var n1 big.Int
	n1.SetUint64(num)
	// != 0 check is done inside the Rsh function
	// Rsh: https://golang.org/src/math/big/int.go?s=25450:25488#L990
	// shr: https://golang.org/src/math/big/nat.go#L1006
	return n1.Rsh(&n1, n)
}

// GetFileIdentifierByPath gets a fileidentifier by path
// it just opens the path and calls GetFileIdentifierByFile
func GetFileIdentifierByPath(path string) (*FileIdentifier, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("GetFileIdentifierByPath open path %v error: %v", path, err)
	}
	return GetFileIdentifierByFile(f)
}

// GetGlobalFileID returns the file id
func GetGlobalFileID(i os.FileInfo) *big.Int {
	return GetFileIdentifier(i).GetGlobalFileID()
}

// for tests
func iterateAllUint64(max uint64, cb func(count uint64)) {
	i := uint64(0)
	addvalue := uint64(0)
	lastValie := i
	if max == 0 {
		max = addvalue - 1
	}
	for ; i >= lastValie && i < max; i = i + addvalue {
		cb(i)
		addvalue = addvalue*3000/2 + 100
		lastValie = i
	}
	cb(max)
}
