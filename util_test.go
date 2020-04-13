package fileidentifier

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
