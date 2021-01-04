package mcts

import "me.dev/go-board-game/common"

// Propagate update Node tree with new score found from playout
func Propagate(node Node, score float32, config Config) {
	move := node.GetMove()
	playoutPlayer := move.GetPlayer()
	count := 0
	for currentNode := node; ; currentNode = currentNode.GetParent() {
		config.PropagateToNode(playoutPlayer, currentNode, score, count)
		if currentNode.IsRoot() {
			break
		}
		count++
	}
}

// BasicPropgateUpdate update the score on a particular node.  Be sure to flip the score when node is connected to a different player
func BasicPropgateUpdate(playoutPlayer common.Player, node Node, score float32, count int) {
	if playoutPlayer != node.GetMove().GetPlayer() && !node.IsRoot() {
		score = 1.0 - score
	}

	node.addScore(score)
}
