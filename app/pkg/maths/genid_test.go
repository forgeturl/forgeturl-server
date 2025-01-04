package maths

import (
	"math"
	"testing"
)

func TestBase58Encode(t *testing.T) {
	str := Base58Encode(math.MaxInt64) // NQm6nKp8qFC
	// str = Base58Encode(math.MaxInt32)  // 4GmR58
	t.Log(str) // NQm6nKp8qFC
}
