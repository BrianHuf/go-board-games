package mcts

// SelectByMostVisits ...
func SelectByMostVisits(root Node) (best Node) {
	maxVisits := 0
	for _, node := range *root.GetChildren() {
		if node.GetVisits() > maxVisits {
			best = node
			maxVisits = node.GetVisits()
		}
	}
	return
}
