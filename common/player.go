package common

// GamePlayer is always used as the first move in a Game
type GamePlayer struct {
}

// ToString ...
func (a GamePlayer) String() string {
	return "<opening>"
}

// Equals ...
func (a GamePlayer) Equals(b Player) bool {
	_, ok := b.(GamePlayer)
	return ok
}

// TwoPlayer common class to suport two player games
type TwoPlayer struct {
	IsPlayer1 bool
}

// ToString ...
func (a TwoPlayer) String() string {
	if a.IsPlayer1 {
		return "p1"
	}
	return "p2"
}

// Equals ...
func (a TwoPlayer) Equals(b Player) bool {
	other, ok := b.(TwoPlayer)
	return ok && (a.IsPlayer1 == other.IsPlayer1)
}
