package mcts

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"me.dev/go-board-game/common"
	"me.dev/go-board-game/games/tictactoe"
)


func Test_FindBestMove_OneThread(t *testing.T) {
	emptyBoard := tictactoe.NewGame()
	bestNode, _ := FindBestMove(emptyBoard, BasicConfig())
	bestMove := bestNode.GetMove()
	assert.True(t, tictactoe.IsCenter(bestMove))
}

func Test_FindBestMove_MultiThreadedSpeedup(t *testing.T) {
	ITER := 10000
	MULTITHREAD := 20

	emptyBoard := tictactoe.NewGame()

	startSingleThreaded := time.Now()
	bestNodeSingleThreaded, rootNodeSingleThreaded := FindBestMove(emptyBoard, slowBasicConfig(ITER))
	bestMoveSingleThreaded := bestNodeSingleThreaded.GetMove()
	elapsedSingleThreaded := time.Since(startSingleThreaded)
	assert.EqualValues(t, ITER, rootNodeSingleThreaded.GetVisits())
	assert.True(t, tictactoe.IsCenter(bestMoveSingleThreaded))

	startMultiThreaded := time.Now()
	bestNodeMultiThreaded, rootNodeMultiThreaded := FindBestMove(emptyBoard, slowMultithreadedConfig(ITER, MULTITHREAD))
	bestMoveMultiThreaded := bestNodeMultiThreaded.GetMove()
	elapsedMultiThreaded := time.Since(startMultiThreaded)
	assert.EqualValues(t, ITER, rootNodeMultiThreaded.GetVisits())
	assert.True(t, tictactoe.IsCenter(bestMoveMultiThreaded))

	speedup := elapsedSingleThreaded.Microseconds() / elapsedMultiThreaded.Microseconds()
	fmt.Printf("speedup %d vs %d\n", speedup, MULTITHREAD)
	assert.Greater(t, int(speedup), MULTITHREAD)
}


// FIXTURES -- slow down the playout to give parallel processing a chance to overcome
// overhead associated with channels
func slowPlayoutSelector(m common.Move) (nextMove common.Move) {
	time.Sleep(10 * time.Microsecond)
	return BasicPlayoutSelector(m)
}

func slowBasicConfig(iters int) Config {
	return Config{
		TreeBuilder:     FixedIterations(iters),
		TreeSelector:    SelectByMostVisits,
		PlayoutSelector: slowPlayoutSelector,
		PlayoutScorer:   BasicPlayoutScorer,
		PlayoutMaxMoves: 10,
		ExpandSelector:  BasicExpandSelect,
		SelectorWeigher: BasicWeigher,
		PropagateToNode: BasicPropgateUpdate}
}

func slowMultithreadedConfig(iters int, threads int) Config {
	return Config{
		TreeBuilder:     FixedIterationsParallel(iters, threads),
		TreeSelector:    SelectByMostVisits,
		PlayoutSelector: slowPlayoutSelector,
		PlayoutScorer:   BasicPlayoutScorer,
		PlayoutMaxMoves: 10,
		ExpandSelector:  BasicExpandSelect,
		SelectorWeigher: BasicWeigher,
		PropagateToNode: BasicPropgateUpdate}
}

