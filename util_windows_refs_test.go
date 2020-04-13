// +build refs

package fileidentifier

import (
	"math/big"
)

// iterateAllFileIdentifier is for for tests
func iterateAllFileIdentifier(cb func(globalId, expectedFileID *big.Int, vol, idxHi, idxLo uint64, f FileIdentifier)) {
	f := FileIdentifier{}

	expected := big.NewInt(0)
	iterateAllUint64(18446744073709551615, func(vol uint64) {
		bigIdxVol := getBigInt(vol, 128)
		expected.Add(expected, bigIdxVol)
		iterateAllUint64(18446744073709551615, func(idxHi uint64) {
			bigIdxHi := getBigInt(idxHi, 64)
			expected.Add(expected, bigIdxHi)
			iterateAllUint64(18446744073709551615, func(idxLo uint64) {
				f.vol = vol
				f.idxHi = idxHi
				f.idxLo = idxLo

				bigIdxLo := getBigInt(idxLo, 0)
				expected.Add(expected, bigIdxLo)

				expectedFileID := big.NewInt(0)
				expectedFileID.Add(bigIdxHi, bigIdxLo)

				cb(expected, expectedFileID, vol, idxHi, idxLo, f)

				expected.Sub(expected, bigIdxLo)
			})
			expected.Sub(expected, bigIdxHi)
		})
		expected.Sub(expected, bigIdxVol)
	})
}
