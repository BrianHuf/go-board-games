package mcts

import (
	"fmt"
	"strings"

	"me.dev/go-board-game/common"
)

// Node is a particular move in a game that is connected in a MSCTS tree
type Node interface {
	IsRoot() bool
	GetParent() Node
	HasChildren() bool
	GetChildren() *[]Node
	setChildren(children *[]common.Move)
	addScore(score float32)
	GetMove() common.Move
	GetScore() float32
	GetVisits() int
}

// NodeSingleThread ...
type NodeSingleThread struct {
	move     common.Move
	parent   *Node
	children []Node
	visits   int
	score    float32
}

func newSingleThreadedRootNode(move common.Move) Node {
	return &NodeSingleThread{move: move}
}

// GetParent ...
func (node *NodeSingleThread) GetParent() Node {
	return *node.parent
}

// IsRoot ...
func (node *NodeSingleThread) IsRoot() bool {
	return node.parent == nil
}

func (node *NodeSingleThread) setChildren(moves *[]common.Move) {
	children := make([]Node, len(*moves))

	var iNode Node = node
	var pNode *Node = &iNode
	for i, move := range *moves {
		nst := &NodeSingleThread{move: move, parent: pNode}
		children[i] = nst
	}

	node.children = children
}

// GetMove ..
func (node *NodeSingleThread) GetMove() common.Move {
	return node.move
}

// HasChildren True when node has already been expanded to have children
func (node *NodeSingleThread) HasChildren() bool {
	return len(node.children) > 0
}

// GetScore 0.0 means this Node is a bad move.  1.0 means this is a good move
func (node *NodeSingleThread) GetScore() float32 {
	return node.score
}

// GetVisits ...
func (node *NodeSingleThread) GetVisits() int {
	return node.visits
}

// GetChildren ...
func (node *NodeSingleThread) GetChildren() (children *[]Node) {
	return &node.children
}

func (node *NodeSingleThread) addScore(score float32) {
	node.score += score
	node.visits++
}

func (node *NodeSingleThread) String() string {
	return node.paddedString(0, 1)
}

func (node *NodeSingleThread) paddedString(depth int, maxDepth int) string {
	node.move.MoveString()
	var ret = fmt.Sprintf("%s%s (%3.2f/%d)\n", 
		strings.Repeat("  ", depth),
		node.move.MoveString(),
		node.GetScore(),
		node.GetVisits())

	if depth < maxDepth {
		for _, child := range node.children {
			ret += child.(*NodeSingleThread).paddedString(depth+1, maxDepth)
		}	
	}

	return ret
}