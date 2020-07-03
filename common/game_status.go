package common

// GameStatusData simple data structure
type GameStatusData struct {
	isDone bool
	winner Player
}

// NewGameStatusWinner ...
func NewGameStatusWinner(p Player) GameStatus {
	return GameStatusData{isDone: true, winner: p}
}

// NewGameStatusTied ...
func NewGameStatusTied() GameStatus {
	return GameStatusData{isDone: true, winner: nil}
}

// NewGameStatusInProgress ...
func NewGameStatusInProgress() GameStatus {
	return GameStatusData{isDone: false, winner: nil}
}

// IsDone ...
func (status GameStatusData) IsDone() bool {
	return status.isDone
}

// GetWinner ...
func (status GameStatusData) GetWinner() Player {
	return status.winner
}

func (status GameStatusData) String() string {
	if status.isDone {
		if status.winner == nil {
			return "tie"
		}
		return "winner " + status.winner.String()
	}
	return "in progress"
}
