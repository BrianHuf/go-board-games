package common

// GameStatusData simple data structure
type GameStatusData struct {
	Done   bool   `json:"isDone"`
	Winner Player `json:"winner"`
}

// NewGameStatusWinner ...
func NewGameStatusWinner(p Player) GameStatus {
	return &GameStatusData{Done: true, Winner: p}
}

// NewGameStatusTied ...
func NewGameStatusTied() GameStatus {
	return &GameStatusData{Done: true, Winner: PlayerNoOne}
}

// NewGameStatusInProgress ...
func NewGameStatusInProgress() GameStatus {
	return &GameStatusData{Done: false, Winner: PlayerNoOne}
}

// IsDone ...
func (status *GameStatusData) IsDone() bool {
	return status.Done
}

// GetWinner ...
func (status *GameStatusData) GetWinner() Player {
	return status.Winner
}

func (status *GameStatusData) String() string {
	if status.Done {
		if status.Winner == PlayerNoOne {
			return "tie"
		}
		return "winner " + status.Winner.String()
	}
	return "in progress"
}
