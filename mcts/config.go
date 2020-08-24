package mcts

import "me.dev/go-board-game/common"

// TreeBuilderType give a root node with coordinate adding children (e.g. time-based, fixed iterations, multi-threaded, etc...)
type TreeBuilderType func(root Node, config Config)

// TreeSelectorType select the best node/move given a fully populated Node tree
type TreeSelectorType func(root Node) (best Node)

// SelectorWeigherType interface to weigh nodes (which child to pick?)
type SelectorWeigherType func(node Node) (weight float32)

// ExpandSelectorType interface to weigh nodes (which new child node to select?)
type ExpandSelectorType func(nodes *[]Node) (selected Node)

// PlayoutSelectorType function to select the next move (which available move to pick when playing out?)
type PlayoutSelectorType func(currentMove common.Move) (nextMove common.Move)

// PlayoutScorerType function to calculate a playout score (0.0 bad - 1.0 good) (is the current node/move any good?)
type PlayoutScorerType func(finalMove common.Move, numMoves int, startingPlayer common.Player, winningPlayer common.Player) (score float32)

// PropgateUpdateType function to update Node scores based on a playout
type PropgateUpdateType func(playoutPlayer common.Player, node Node, score float32, count int)

// Config control some key areas of the MCTS algorithm
type Config struct {
	TreeBuilder     TreeBuilderType
	TreeSelector    TreeSelectorType
	PlayoutSelector PlayoutSelectorType
	PlayoutScorer   PlayoutScorerType
	PlayoutMaxMoves int
	ExpandSelector  ExpandSelectorType
	SelectorWeigher SelectorWeigherType
	PropagateToNode PropgateUpdateType
}

// BasicConfig standard MCTS behavior
func BasicConfig() Config {
	return Config{
		TreeBuilder:     FixedIterations(10000),
		TreeSelector:    SelectByMostVisits,
		PlayoutSelector: BasicPlayoutSelector,
		PlayoutScorer:   BasicPlayoutScorer,
		PlayoutMaxMoves: 10,
		ExpandSelector:  BasicExpandSelect,
		SelectorWeigher: BasicWeigher,
		PropagateToNode: BasicPropgateUpdate}
}
