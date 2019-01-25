package skiplist

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"testing"
)

func BenchmarkRNG(b *testing.B) {
	max := big.NewInt(math.MaxInt64)

	for i := 0; i < b.N; i++ {
		_, err := crand.Int(crand.Reader, max)
		if err != nil {
			b.Fatalf("calling rand.Int: %s", err)
		}
	}
}

func BenchmarkPRNG(b *testing.B) {
	randSrc := rand.NewSource(1234)
	randGen := rand.New(randSrc)

	for i := 0; i < b.N; i++ {
		res := randGen.Int63()
		if res == 0 {
			b.Fatalf("dummy check for parity with BenchmarkRNG")
		}
	}
}
