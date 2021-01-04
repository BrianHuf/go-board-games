package tictactoe

import (
	"fmt"
	"math/bits"
	"strconv"

	"me.dev/go-board-game/common"
)

/*
  1 |  2|  4
  --|---|---
  8 | 16| 32
  --|---|---
  64|128|256
*/

var winningBoardMasks []uint32 = []uint32{
	7, 56, 448, // horizontal
	73, 146, 292, // vertical
	84, 273} // diagonal

// TicTacToeMove ...
type TicTacToeMove struct {
	state    uint32
	previous *TicTacToeMove
}

// StateJSON ...
type StateJSON struct {
	State      string `json:"state"`
	LastMove   string `json:"lastMove"`
	IsDone     bool   `json:"isDone"`
	Winner     string `json:"winner,omitempty"`
	NextPlayer string `json:"nextPlayer,omitempty"`
}

// GetJSON ...
func (a *TicTacToeMove) GetJSON() interface{} {
	status := a.GetGameStatus()

	var winner string = ""
	var nextPlayer string = ""
	if status.IsDone() {
		if status.GetWinner() != common.PlayerNoOne {
			winner = status.GetWinner().String()
		}
	} else {
		nextPlayer = a.nextPlayer().String()
	}

	return StateJSON{
		State:      a.getBoardString(),
		LastMove:   a.MoveString(),
		IsDone:     status.IsDone(),
		Winner:     winner,
		NextPlayer: nextPlayer}
}

// GetPlayer ...
func (a *TicTacToeMove) GetPlayer() common.Player {
	if a.previous == nil {
		return common.PlayerNoOne
	}
	if isPlayer1(a.state) {
		return common.Player1
	}
	return common.Player2
}

func (a *TicTacToeMove) nextPlayer() common.Player {
	if a.previous == nil || !isPlayer1(a.state) {
		return common.Player1
	} else {
		return common.Player2
	}
}

// GetPrevious ...
func (a *TicTacToeMove) GetPrevious() common.Move {
	if a.previous == nil {
		return nil
	}
	return (*a).previous
}

// String ...
func (a *TicTacToeMove) String() string {
	return "board"
}

// NextAvailableMoves ...
func (a *TicTacToeMove) NextAvailableMoves() (available []common.Move) {
	isNextPlayerP1 := !isPlayer1(a.state)
	board := getCombinedBoard(a.state)
	for i := 0; i < 9; i++ {
		bit := uint32(1) << i
		if board&bit == 0 {
			if !isNextPlayerP1 {
				bit = bit << 9
			}

			newState := flipPlayer(a.state) | bit
			m := TicTacToeMove{previous: a, state: newState}
			available = append(available, &m)
		}
	}

	return
}

// GetGameStatus ...
func (a *TicTacToeMove) GetGameStatus() common.GameStatus {
	if a.previous == nil {
		return common.NewGameStatusInProgress()
	}

	if getCombinedBoard(a.state) == 511 {
		return common.NewGameStatusTied()
	}

	p1 := isPlayer1(a.state)
	board := getBoard(a.state, p1)
	if isWinningBoard(board) {
		return common.NewGameStatusWinner(a.GetPlayer())
	}

	return common.NewGameStatusInProgress()
}

// BoardString CLI based represention of game give current move
func (a *TicTacToeMove) BoardString() string {
	b := a.getBoardString()

	return fmt.Sprintf("%c|%c|%c\n-----\n%c|%c|%c\n-----\n%c|%c|%c\n",
		b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7], b[8])
}

func (a *TicTacToeMove) getBoardString() string {
	board1 := getBoard(a.state, true)
	board2 := getBoard(a.state, false)

	var b = ""
	for i := 0; i < 9; i++ {
		bit := uint32(1) << i
		if board1&bit > 0 {
			b += "X"
		} else if board2&bit > 0 {
			b += "O"
		} else {
			b += " "
		}
	}

	return b
}

// MoveString represents the current move
func (a *TicTacToeMove) MoveString() string {
	if a.previous == nil {
		return "-"
	}

	diff := (a.state - (*(a.previous)).state) & 262143
	if !isPlayer1(a.state) {
		diff = diff >> 9
	}
	playedIndex := bits.TrailingZeros32(diff)
	return strconv.Itoa(playedIndex)
}

// PlayMovesByIndex ...
func (a *TicTacToeMove) PlayMovesByIndex(moves *[]int) common.Move {
	return common.PlayMovesByIndex(a, moves)
}

// PlayMovesByString ...
func (a *TicTacToeMove) PlayMovesByString(moves string) common.Move {
	return common.PlayMovesByString(a, moves, 1)
}

// NewTicTacToeGame start a new TicToe game by returning the opening move
func NewTicTacToeGame() common.Move {
	return &TicTacToeMove{262144, nil}
}

// IsCenter return True is move is the in the center of the board
func IsCenter(move common.Move) bool {
	return move.MoveString() == "4"
}

// local functions
func isPlayer1(state uint32) bool {
	return state&262144 == 0
}

func flipPlayer(state uint32) uint32 {
	return state ^ 262144
}

func getBoard(state uint32, isPlayer1 bool) uint32 {
	if isPlayer1 {
		return state & 511
	}

	return (state & 261632) >> 9
}

func getCombinedBoard(state uint32) uint32 {
	return getBoard(state, true) | getBoard(state, false)
}

func isWinningBoard(board uint32) bool {
	for _, check := range winningBoardMasks {
		ucheck := uint32(check)
		if ucheck&board == ucheck {
			return true
		}
	}
	return false
}
