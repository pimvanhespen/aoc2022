package main

import "testing"

var testData = `GbccTtTSGGbgrcWBGGrdgTnVQnCmNpCJlNnNPVfClcnN
vMzvZhzhwDLVmQnClwwNQp
FRsZFzjQFsqRzRRjDZbdtTgdHBBWGrdBdHHs
GbccTtTSGGbgrcWBGGrdgTnVQnCmNpCJlNnNPVfClcnN
vMzvZhzhwDLVmQnClwwNQp
FRsZFzjQFsqRzRRjDZbdtTgdHBBWGrdBdHHs
`

func Test_solve2_bytes(t *testing.T) {
	v, err := solve2_bytes([]byte(testData))
	if err != nil {
		t.Fatal(err)
	}
	if v <= 0 {
		t.Fatal("expected positive value")
	}
}
