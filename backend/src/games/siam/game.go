package siam

import (
	"bytes"

	"me.dev/go-board-game/common"
)

/*
  SIAM is played on a 5x5 board
  Each players has 5 directional pieces
  The board have three mountains
  The winning player pushes a mountain off the board
*/

// Players
var (
	PLAYER_1 = common.NewPlayer(0)
	PLAYER_2 = common.NewPlayer(1)
)

// Board
var BOARD = common.Board(5)


// Direction ...
type Direction int
const (
    Up Direction = iota
    Down
    Left
    Right
)

func (d Direction) String() string {
    return [...]string{"Up", "Down", "Left", "Right"}[d]
}

// Piece ...
type Piece byte

// String ...
func (p Piece) String() string {
	return [...]string{"(empty)", "Mountain", 
		"P1 Up", "P1 Down", "P1 Right", "P1 Left",
		"P2 Up", "P2 Down", "P2 Right", "P2 Left"}[p]
}

// Char ...
func (p *Piece) Char() string {
	if p == nil {
		return "."
	}
	return [...]string{".", "M", 
		"U", "D", "R", "L",
		"u", "d", "r", "l"}[*p]
}

var (
	constantPieceEmpty = Piece(0)
	constantPieceMountain = Piece(1)
)

// PiecesByPlayer ...
func PiecesByPlayer(p common.Player) []Piece {
	if p.Equals(PLAYER_1) {
		return []Piece{
			Piece(2), Piece(3), Piece(4), Piece(5)}
	} else if p.Equals(PLAYER_2) {
		return []Piece{
			Piece(6), Piece(7), Piece(8), Piece(9)}
	}

	return nil
}

// FromPlayer ...
func (p Piece) FromPlayer(player common.Player) bool {
	return (player.Equals(PLAYER_1) && byte(p) > 1 && byte(p) < 6) ||
		(player.Equals(PLAYER_2) && byte(p) > 5)
}

// Move ...
type Move struct {
	playedBy common.Player
	previous *Move
	board [25]*Piece
	from *common.Location
	to *common.Location
}

// StateJSON ...
type StateJSON struct {
	State string `json:"state"`
	LastMove string `json:"lastMove"`
	IsDone bool `json:"isDone"`
	Winner string `json:"winner,omitempty"`
	NextPlayer string `json:"nextPlayer,omitempty"`
}

// GetJSON ...
func (a *Move) GetJSON() interface{} {
	status := a.GetGameStatus()

	var winner string = ""
	var nextPlayer string = ""
	if status.IsDone() {
		if (status.GetWinner() != nil) {
			winner = status.GetWinner().String()
		}
	} else {
		nextPlayer = a.nextPlayer().String()
	}

	return StateJSON{
		State: a.BoardString(), 
		LastMove: a.MoveString(),
		IsDone: status.IsDone(),
		Winner: winner,
		NextPlayer: nextPlayer}
}

// GetPlayer ...
func (a *Move) GetPlayer() common.Player {
	return a.playedBy
}

func (a *Move) nextPlayer() common.Player {
	if a.GetPlayer() == PLAYER_1 {
		return PLAYER_2
	}

	return PLAYER_1
}

// GetPrevious ...
func (a *Move) GetPrevious() common.Move {
	if a.previous == nil {
		return nil
	}
	return (*a).previous
}

// String ...
func (a *Move) String() string {
	return "board"
}

// NextAvailableMoves ...
func (a *Move) NextAvailableMoves() (available []common.Move) {
	nextPlayer := a.nextPlayer()
	freePieces := 5

	var addPieceMoves []common.Move

	for index, piece := range(a.board) {
		location := common.Location(byte(index))

		if BOARD.IsEdge(location) {
			if piece == nil {
				// ADD piece
				for _, newPiece := range(PiecesByPlayer(nextPlayer)) {
					var newBoard = a.board
					var newPieceCopy = newPiece
					newBoard[index] = &newPieceCopy
	
					var newMove = Move{
						playedBy: nextPlayer,
						previous: a,
						board: newBoard,
						from: nil,
						to: &location}

					addPieceMoves = append(addPieceMoves, &newMove)
				}
			} else if piece.FromPlayer(nextPlayer) {
				// REMOVE piece
				newBoard := a.board
				newBoard[index] = nil

				newMove := Move{
					playedBy: nextPlayer,
					previous: a,
					board: newBoard,
					from: &location,
					to: nil}

				available = append(available, &newMove)
			}
		}

		if piece != nil {			
			freePieces--

			// MOVE piece (TODO)			
		}
	} 

	if freePieces > 0 {
		available = append(available, addPieceMoves...)
	}

	return
}

// GetGameStatus ...
func (a *Move) GetGameStatus() common.GameStatus {
	var mountainCount int = 0
	for _, piece := range(a.board) {
		if piece == &constantPieceMountain {
			mountainCount++
		}
	}

	if mountainCount == 3 {
		return common.NewGameStatusInProgress()
	}

	return common.NewGameStatusInProgress()
}

// BoardString CLI based represention of game give current move
func (a *Move) BoardString() string {
	var b bytes.Buffer
	for _, p := range(a.board) {
		if p == nil {
			b.WriteString(constantPieceEmpty.Char())
		} else {
			b.WriteString(p.Char())
		}
	}

	return b.String()
}


// MoveString represents the current move
func (a *Move) MoveString() string {
	return a.from.String() + a.to.String() + a.toPiece().Char()
}

func (a *Move) toPiece() *Piece {
	if a.to == nil {
		return nil
	}

	index := int(*a.to)
	return a.board[index]
}

// PlayMovesByIndex ...
func (a *Move) PlayMovesByIndex(moves *[]int) common.Move {
	return common.PlayMovesByIndex(a, moves)
}

// PlayMovesByString ...
func (a *Move) PlayMovesByString(moves string) common.Move {
	return common.PlayMovesByString(a, moves, 3)
}

// NewGame start a new TicToe game by returning the opening move
func NewGame() common.Move {
	initialBoard := [25]*Piece{} 
	initialBoard[11] = &constantPieceMountain
	initialBoard[12] = &constantPieceMountain
	initialBoard[13] = &constantPieceMountain

	return &Move{
		playedBy: common.NoPlayer{},
		previous: nil,
		board: initialBoard,
		from: nil,
		to: nil}
}
