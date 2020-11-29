package mcts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"me.dev/go-board-game/games/tictactoe"
)

func Test_FindBestMove(t *testing.T) {
	emptyBoard := tictactoe.NewGame()
	bestMove, _ := FindBestMove(emptyBoard, BasicConfig())
	assert.Equal(t, "4", bestMove.GetMove().MoveString())
}
