package mcts

import "math/rand"

// ExpandSelector interface to weigh nodes
type ExpandSelector func(nodes []Node) (selected Node)

// BasicExpand get next
func BasicExpand(selected Node) (expanded Node) {
	return Expand(selected, BasicExpandSelect)
}

// Expand ...
func Expand(selected Node, selector ExpandSelector) (expanded Node) {
	if selected.visits == 0 || selected.move.GetGameStatus().IsDone() {
		return selected
	}

	for _, move := range selected.move.NextAvailableMoves() {
		selected.addChild(move)
	}

	expanded = selector(selected.children)
	return
}

// BasicExpandSelect random uniform selection
func BasicExpandSelect(nodes []Node) (selected Node) {
	len := len(nodes)
	selected = nodes[rand.Intn(len)]
	return
}
