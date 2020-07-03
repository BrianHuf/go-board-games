package common

import "strconv"

// NoPlayer is always used as the first move in a Game
type NoPlayer struct {
}

// ToString ...
func (a NoPlayer) String() string {
	return "<opening>"
}

// Equals ...
func (a NoPlayer) Equals(b Player) bool {
	_, ok := b.(NoPlayer)
	return ok
}

// NewPlayer ...
func NewPlayer(num uint) PlayerData {
	return PlayerData{num}
}

// PlayerData common class to suport n-player games
type PlayerData struct {
	num uint
}

// ToString ...
func (a PlayerData) String() string {
	return "p" + strconv.Itoa(int(a.num+1))
}

// Equals ...
func (a PlayerData) Equals(b Player) bool {
	other, ok := b.(PlayerData)
	return ok && (a.num == other.num)
}
