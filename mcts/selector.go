package mcts

import (
	"math"
	"math/rand"
)

// Selector select a node to expand
func Selector(root Node, config Config) Node {
	for current := root; ; current = pickNext(current, config.SelectorWeigher) {
		if !current.HasChildren() {
			return current
		}
	}
}

func pickNext(node Node, weigher SelectorWeigherType) (selected Node) {
	children := *node.GetChildren()

	var weights []float32 = make([]float32, len(children))
	for i, node := range children {
		weights[i] = weigher(node)
	}

	index := weightedRandom(weights)
	selected = children[index]
	return
}

var c float32 = float32(math.Sqrt(2))

// BasicWeigher textbook mcts node selection
// FIXME equation adjust to avoid NaN and inf
func BasicWeigher(node Node) (weight float32) {
	wi := node.GetScore()
	ni := float32(node.GetVisits()) + 1.0
	Ni := float64(node.GetParent().GetVisits()) + 1.0
	weight = float32(wi/ni + c*float32(math.Log(Ni))/ni)
	if weight < 0 {
		return 0
	}
	return
}

func weightedRandom(weights []float32) int {
	// https://zliu.org/post/weighted-random/
	n := len(weights)
	if n == 0 {
		return 0
	}
	cdf := make([]float32, n)
	var sum float32 = 0.0
	for i, w := range weights {
		if i > 0 {
			cdf[i] = cdf[i-1] + w
		} else {
			cdf[i] = w
		}
		sum += w
	}
	r := rand.Float32() * sum
	var l, h int = 0, n - 1
	for l <= h {
		m := l + (h-l)/2
		if r <= cdf[m] {
			if m == 0 || (m > 0 && r > cdf[m-1]) {
				return m
			}
			h = m - 1
		} else {
			l = m + 1
		}
	}
	return -1
}
