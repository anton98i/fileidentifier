// +build !refs

package fileidentifier

import (
	"math/big"
)

func iterateAllFileIdentifier(cb func(globalId, expectedFileID *big.Int, vol, idxHi, idxLo uint64, f FileIdentifier)) {
	f := FileIdentifier{}

	expected := big.NewInt(0)
	iterateAllUint64(4294967295, func(vol uint64) {
		bigIdxVol := getBigInt(vol, 64)
		expected.Add(expected, bigIdxVol)
		iterateAllUint64(4294967295, func(idxHi uint64) {
			bigIdxHi := getBigInt(idxHi, 32)
			expected.Add(expected, bigIdxHi)
			iterateAllUint64(4294967295, func(idxLo uint64) {
				f.vol = uint32(vol)
				f.idxHi = uint32(idxHi)
				f.idxLo = uint32(idxLo)

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
