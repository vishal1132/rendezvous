package rendezvous

import "testing"

var r = Rendezvous{
	nodes: []string{"a", "b", "c", "d", "e", "f", "g", "h"},
	f:     HashFunction(MD5),
}

func BenchmarkGetScore(b *testing.B) {

	for i := 0; i < b.N; i++ {
		r.GetScore([]byte("abcd"))
	}
}
