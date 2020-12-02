package server

// GameDto ...
type GameDto struct {
	State interface{} `json:"state"`
	Moves []MoveDto `json:"moves"`
}

// MoveDto ...
type MoveDto struct {
	Value interface{} `json:"value"`
}
