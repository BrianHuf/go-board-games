package common

import "fmt"

// Player is a player in the game
type Player interface {
	Equals(other Player) bool
	fmt.Stringer
}

// GameStatus interface
type GameStatus interface {
	IsDone() bool
	GetWinner() Player
	fmt.Stringer
}

// Move is a particular move in a game played by a players
// The game begins with a new Game Move which is associated
// with the GamePlayer
type Move interface {
	BoardString() string
	GetPlayer() Player
	GetPrevious() Move
	GetGameStatus() GameStatus
	MoveString() string
	NextAvailableMoves() []Move
	PlayMovesByIndex(moves []int) Move
	PlayMovesByString(moves string) Move
	fmt.Stringer
}
