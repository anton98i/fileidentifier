package fileidentifier

import "testing"

func expect(t *testing.T, excpectedVal, isVal interface{}) {
	if excpectedVal != isVal {
		t.Error("expect failed, expected:", excpectedVal, "received:", isVal)
	}
}

func expectTrue(t *testing.T, val bool) {
	expect(t, true, val)
}

func expectNil(t *testing.T, val interface{}) {
	expect(t, nil, val)
}
