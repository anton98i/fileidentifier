package fileidentifier

import (
	"math/big"
)

// FileIdentifier interface
type FileIdentifier interface {
	// GetGlobalFileID returns the device id + file id combined to one id (a "uint128")
	GetGlobalFileID() *big.Int

	// GetDeviceID returns the device id (on windows it is a uint32 casted as uint64)
	GetDeviceID() uint64

	// GetFileID returns the file id
	GetFileID() uint64
}

// FileIdentEx interface
type FileIdentEx interface {
	// GetGlobalFileID returns the device id + file id combined to one id (a "uint192")
	GetGlobalFileID() *big.Int

	// GetDeviceID returns the device id
	GetDeviceID() uint64

	// GetFileID returns the file id as a "uint128"
	GetFileID() *big.Int
}

// getBigInt returns num << n
func getBigInt(num uint64, n uint) *big.Int {
	var n1 big.Int
	n1.SetUint64(num)
	// != 0 check is done inside the shl function
	// Lsh: https://golang.org/src/math/big/int.go?s=25314:25352#L993
	// shl: https://golang.org/src/math/big/nat.go#L981
	return n1.Lsh(&n1, n)
}
