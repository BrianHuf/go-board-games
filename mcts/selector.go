package mcts

import (
	"math"
	"math/rand"
)

// SelectorWeigher interface to weigh nodes
type SelectorWeigher func(node Node) (weight float32)

// BasicSelector repeatedly pick node children until more children are found
func BasicSelector(root Node) Node {
	for current := root; ; current = pickNext(current, BasicWeight) {
		if !current.hasChildren() {
			return current
		}
	}
}

func pickNext(node Node, weigher SelectorWeigher) (selected Node) {
	var weights []float32 = make([]float32, len(node.children))
	for i, node := range node.children {
		weights[i] = weigher(node)
	}
	index := weightedRandom(weights)
	selected = node.children[index]
	return
}

var c float64 = math.Sqrt(2)

// BasicWeight textbook mcts node selection
func BasicWeight(node Node) (weight float32) {
	wi := node.score
	ni := float64(node.visits)
	Ni := float64(node.parent.visits) + 1.0
	weight = float32(wi/ni + c*math.Log(Ni)/ni)
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
