package mcts

import (
	"fmt"
	"testing"

	"me.dev/go-board-game/games/tictactoe"
)

func Test_FindBestMove(t *testing.T) {
	move := tictactoe.NewGame()
	bestMove := FindBestMove(move, BasicConfig())
	fmt.Println("Best move = ", bestMove.GetMove().BoardString())
}
