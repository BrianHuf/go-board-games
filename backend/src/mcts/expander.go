package mcts

import "math/rand"

// Expand ...
func Expand(selected Node, config Config) (newNode Node) {
	if shouldExpand(selected) {
		addAvailableMoves(selected)
		children := selected.GetChildren()
		newNode = config.ExpandSelector(children)
		return
	}
	return selected
}

func shouldExpand(node Node) bool {
	return (node.IsRoot() || node.GetVisits() == 0) && !node.GetMove().GetGameStatus().IsDone()
}

func addAvailableMoves(node Node) {
	available := node.GetMove().NextAvailableMoves()
	node.setChildren(&available)
}

// BasicExpandSelect random uniform selection
func BasicExpandSelect(nodes *[]Node) (selected Node) {
	len := len(*nodes)
	selected = (*nodes)[rand.Intn(len)]
	return
}
