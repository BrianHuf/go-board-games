package mcts

import (
	"me.dev/go-board-game/common"
)

// FindBestMove entry point for a MCTS evaluated.  The config argument
// controls key areas of the algorithm
//
// 1.  TreeBuilder is given a root node.  This loops and calls DoOneRound
// 2.    DoOneRound
// 3.      Selector finds a node to expand
// 4.      Expand creates a new child node
// 5.      Playout determines how good the expanded move is
// 6.      Propagator update the score and visits in the Node tree
// 7.  TreeSelect finds the best node
func FindBestMove(move common.Move, config Config) (best Node) {
	root := newSingleThreadedRootNode(move)
	config.TreeBuilder(root, config)
	best = config.TreeSelector(root)
	return
}

// DoRound select, expand, playout, and propagate
func DoRound(root Node, config Config) {
	selected := Selector(root, config)
	new := Expand(selected, config)
	score, _ := Playout(new.GetMove(), config)
	Propagate(new, score, config)
}
