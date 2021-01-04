package games

import (
	"me.dev/go-board-game/mcts"

	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Basic(t *testing.T) {
	game := NewGame()
	assert.Equal(t, "<opening>", game.GetPlayer().String())
	assert.Equal(t, false, game.GetGameStatus().IsDone())

	game.NextAvailableMoves()

	// available1 := game.NextAvailableMoves()
	// assert.Equal(t, 64, len(available1))

	// move1 := available1[0]
	// assert.Equal(t, "p1", move1.GetPlayer().String())
	// assert.Equal(t, false, move1.GetGameStatus().IsDone())

	// available2 := move1.NextAvailableMoves()
	// assert.Equal(t, 60, len(available2))

	// move2 := available2[0]
	// assert.Equal(t, "p2", move2.GetPlayer().String())
	// assert.Equal(t, false, move2.GetGameStatus().IsDone())

	// available3 := move2.NextAvailableMoves()
	// assert.Equal(t, 57, len(available3))

	// json1, _ := json.MarshalIndent(available1[0].GetJSON(), "", " ")
	// t.Log("\n" + string(json1))

	// json2, _ := json.MarshalIndent(available1[4].GetJSON(), "", " ")
	// t.Log("\n" + string(json2))

	// t.Log("\n" + move2.BoardString())
	// t.Log("\n" + move2.MoveString())
	// t.Log("Passed")
}

// func Test_SimpleGameByIndex(t *testing.T) {
// 	lastMove := NewGame().PlayMovesByIndex(&[]int{0, 0, 1, 0, 2})

// 	t.Log("\n" + lastMove.BoardString())
// 	status := lastMove.GetGameStatus()
// 	isDone := status.IsDone()
// 	assert.Equal(t, true, isDone)
// 	assert.Equal(t, "p1", status.GetWinner().String())
// }

func Test_SimpleGameByString(t *testing.T) {
	game := NewGame()
	lastMove := game.PlayMovesByString("UvUUxuvqUxsuqlURfrlgUfgu")

	t.Log("\n" + lastMove.BoardString())
	assert.False(t, lastMove.GetGameStatus().IsDone())
	assert.Equal(t, 20, len(lastMove.NextAvailableMoves()))
	lastMove.GetGameStatus()
	t.Log("Passed")
}

func Test_Ai(t *testing.T) {
	game := NewGame()
	_, root := mcts.FindBestMove(game, mcts.BasicConfig())
	assert.NotNil(t, root)
}

func Test_AiParallel(t *testing.T) {
	game := NewGame()
	_, root := mcts.FindBestMove(game, mcts.MultithreadedConfig())
	assert.NotNil(t, root)
}
