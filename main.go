package main

import (
	"fmt"

	"me.dev/go-board-game/games/tictactoe"
)

func main() {
	fmt.Println("board game cli")
	tictactoe.NewGame()
}
