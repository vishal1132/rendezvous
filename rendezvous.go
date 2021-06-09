package rendezvous

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"sort"
)

// CalcHash takes a key as a slice of byte(so that you can marshal anything into slice of bytes)
// and still use this. it's not interface{}, because you might be marshaling into slice of bytes by
// your own logic, rather than a single general way to marshal.
type CalcHash func(key []byte, node string) uint64

func summize(b []byte) uint64 {
	var val uint64 = 0
	for _, v := range b {
		val += uint64(v)
	}
	return val
}

type Rendezvous struct {
	// represent nodes as their names/ids.
	nodes []string
	f     CalcHash
}

// Returns a CalcHash type from a list of hashes
func HashFunction(hash HashName) CalcHash {
	var scoreFuncs = map[HashName]CalcHash{
		MD5: func(key []byte, node string) uint64 {
			h := md5.Sum(append(key, []byte(node)...))
			return summize(h[:])
		},
		SHA1: func(key []byte, node string) uint64 {
			h := sha1.Sum(append(key, []byte(node)...))
			return summize(h[:])
		},
		SHA256: func(key []byte, node string) uint64 {
			h := sha256.Sum256(append(key, []byte(node)...))
			return summize(h[:])
		},
	}

	return scoreFuncs[hash]
}

// New creates a new rendezvous type.
// Use this type to call any function further.
// If you want to generate a New rendezvous from a slice of strings, use New(slice...)
func New(f CalcHash, nodes ...string) *Rendezvous {
	return &Rendezvous{
		nodes: nodes,
		f:     f,
	}
}

// AddNodes accepts variadic number of string arguments.
func (r *Rendezvous) AddNodes(nodes ...string) error {
	r.nodes = append(r.nodes, nodes...)
	return nil
}

// GetNTop returns a string slice of size n
func (r *Rendezvous) GetNTop(n int, key []byte) []string {
	if n <= 0 || n >= len(r.nodes) {
		// return all the nodes
		return r.nodes
	}
	scores := make([]uint64, 0, len(r.nodes))
	scoresMap := make(map[uint64]string, len(r.nodes))
	for _, v := range r.nodes {
		val := r.f(key, v)
		scores = append(scores, val)
		scoresMap[val] = v
	}
	// reverse sort
	sort.Slice(scores, func(i, j int) bool { return scores[i] > scores[j] })
	topn := make([]string, 0, n)
	for i := 0; i < n; i++ {
		topn = append(topn, scoresMap[scores[i]])
	}
	return topn
}

// GetScore returns the name of the node with highest score and the socre of that node.
// Highest score is 0 in case of error.
func (n *Rendezvous) GetScore(key []byte) (string, uint64) {
	if len(n.nodes) == 0 {
		return "", 0
	}

	var highestScore uint64 = 0
	var highestNodeID uint64 = 0

	for i, v := range n.nodes {
		val := n.f(key, v)
		if val > highestScore {
			highestScore = val
			highestNodeID = uint64(i)
		}
	}
	return n.nodes[highestNodeID], uint64(highestScore)
}
