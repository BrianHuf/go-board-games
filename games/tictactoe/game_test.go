package tictactoe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Basic(t *testing.T) {
	game := NewGame()
	assert.Equal(t, "<opening>", game.GetPlayer().String())
	assert.Equal(t, nil, game.GetWinner())

	available1 := game.NextAvailableMoves()
	assert.Equal(t, 9, len(available1))

	move1 := available1[0]
	assert.Equal(t, "p1", move1.GetPlayer().String())
	assert.Equal(t, nil, move1.GetWinner())

	available2 := move1.NextAvailableMoves()
	assert.Equal(t, 8, len(available2))

	move2 := available2[0]
	assert.Equal(t, "p2", move2.GetPlayer().String())
	assert.Equal(t, nil, move2.GetWinner())

	t.Log("\n" + move2.BoardString())
	t.Log("Passed")
}

func Test_SimpleGameByIndex(t *testing.T) {
	lastMove := NewGame().PlayMovesByIndex([]int{0, 0, 1, 0, 2})

	t.Log("\n" + lastMove.BoardString())
	assert.Equal(t, "p1", lastMove.GetWinner().String())
	t.Log("Passed")
}

func Test_SimpleGameByString(t *testing.T) {
	lastMove := NewGame().PlayMovesByString("01326")

	t.Log("\n" + lastMove.BoardString())
	assert.Equal(t, "p1", lastMove.GetWinner().String())
	t.Log("Passed")
}
