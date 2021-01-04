package common

import "strconv"

var (
	// PlayerNoOne use for initial game state
	PlayerNoOne = Player(0)

	// Player1 fist player
	Player1 = Player(1)

	// Player2 second player
	Player2 = Player(2)
)

// ToString ...
func (player Player) String() string {
	if player == 0 {
		return "<no-one>"
	}
	return "p" + strconv.Itoa(int(player))
}

// Next P1 -> P2 or P2 -> P1
func (player Player) Next() Player {
	if player == 1 {
		return Player(2)
	} else if player == 2 {
		return Player(1)
	}
	panic("next player not found")
}
