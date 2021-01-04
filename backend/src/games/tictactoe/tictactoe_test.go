package tictactoe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Basic(t *testing.T) {
	game := NewGame()
	assert.Equal(t, "<opening>", game.GetPlayer().String())
	assert.Equal(t, false, game.GetGameStatus().IsDone())

	available1 := game.NextAvailableMoves()
	assert.Equal(t, 9, len(available1))

	move1 := available1[0]
	assert.Equal(t, "p1", move1.GetPlayer().String())
	assert.Equal(t, false, move1.GetGameStatus().IsDone())

	available2 := move1.NextAvailableMoves()
	assert.Equal(t, 8, len(available2))

	move2 := available2[0]
	assert.Equal(t, "p2", move2.GetPlayer().String())
	assert.Equal(t, false, move2.GetGameStatus().IsDone())

	t.Log("\n" + move2.BoardString())
	t.Log("Passed")
}

func Test_SimpleGameByIndex(t *testing.T) {
	lastMove := NewTicTacToeGame().PlayMovesByIndex(&[]int{0, 0, 1, 0, 2})

	t.Log("\n" + lastMove.BoardString())
	status := lastMove.GetGameStatus()
	isDone := status.IsDone()
	assert.Equal(t, true, isDone)
	assert.Equal(t, "p1", status.GetWinner().String())
}

func Test_SimpleGameByString(t *testing.T) {
	lastMove := NewGame().PlayMovesByString("01326")

	t.Log("\n" + lastMove.BoardString())
	assert.Equal(t, "p1", lastMove.GetGameStatus().GetWinner().String())
	t.Log("Passed")
}
