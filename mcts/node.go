package mcts

import (
	"me.dev/go-board-game/common"
)

// Node ...
type Node struct {
	move     common.Move
	parent   *Node
	children []Node
	visits   int
	score    float64
	gameOver bool
}

func newRootNode(move common.Move) Node {
	return Node{move: move}
}

func (node Node) addChild(move common.Move) Node {
	child := Node{move: move, parent: &node}
	node.children = append(node.children, child) // MUTATING/THREADING CONCERN
	return child
}

func (node Node) hasChildren() bool {
	return len(node.children) > 0
}
