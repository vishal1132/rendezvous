package rendezvous

import (
	"crypto/md5"
	"hash"
	"sort"
	"sync"
)

func summize(b []byte) uint64 {
	var val uint64 = 0
	for _, v := range b {
		val += uint64(v)
	}
	return val
}

type rendezvous struct {
	// represent nodes as their names/ids.
	nodes []string
	mu    sync.Mutex
	f     hash.Hash
}

var _ Rendezvous = (*rendezvous)(nil)

type Rendezvous interface {
	GetScore(key []byte) (string, uint64)
	GetNTop(n int, key []byte) []string
}

type nodeScoreSorter []nodescore

func (n nodeScoreSorter) Len() int           { return len(n) }
func (n nodeScoreSorter) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n nodeScoreSorter) Less(i, j int) bool { return n[i].score > n[j].score }

// New creates a new rendezvous type.
// Use this type to call any function further.
// If you want to generate a New rendezvous from a slice of strings, use New(slice...).
func New(f hash.Hash, nodes ...string) Rendezvous {
	r := rendezvous{
		nodes: nodes,
	}
	if f == nil {
		f = md5.New()
	}
	r.f = f
	return &r
}

// AddNodes accepts variadic number of string arguments.
func (r *rendezvous) AddNodes(nodes ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes = append(r.nodes, nodes...)
}

// RemoveNodes removes nodes.
func (r *rendezvous) RemoveNodes(nodes ...string) {
	// remove nodes implementation here.
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(nodes) >= len(r.nodes) {
		r.nodes = make([]string, 0)
		return
	}
	// for _, v := range nodes {
	// }
}

// GetNTop returns a string slice of size n.
func (r *rendezvous) GetNTop(n int, key []byte) []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	if n <= 0 || n >= len(r.nodes) {
		// return all the nodes
		return r.nodes
	}

	scores := make([]nodescore, 0, len(r.nodes))
	r.f.Write(key)

	for _, v := range r.nodes {
		val := summize(r.f.Sum([]byte(v)))
		scores = append(scores, nodescore{node: v, score: val})
	}

	sort.Sort(nodeScoreSorter(scores))
	r.f.Reset()

	var nodes = make([]string, 0, n)

	for i := 0; i < n; i++ {
		nodes = append(nodes, scores[i].node)
	}

	return nodes
}

type nodescore struct {
	node  string
	score uint64
}

// GetScore returns the name of the node with highest score and the socre of that node.
// Highest score is 0 in case of error.
func (r *rendezvous) GetScore(key []byte) (string, uint64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.nodes) == 0 {
		return "", 0
	}

	var highestScore uint64 = 0
	var highestNodeID uint64 = 0

	r.f.Write(key)

	for i, v := range r.nodes {
		val := summize(r.f.Sum([]byte(v)))
		if val > highestScore {
			highestScore = val
			highestNodeID = uint64(i)
		}
	}
	r.f.Reset()

	return r.nodes[highestNodeID], uint64(highestScore)
}
