package tictactoe

import (
	"fmt"
	"math/bits"
	"strconv"

	"me.dev/go-board-game/common"
)

/*
  1 |  2|  4
  8 | 16| 32
  64|128|256
*/

var winningBoardMasks []uint32 = []uint32{
	7, 56, 448, // horizontal
	73, 146, 292, // vertical
	84, 273} // diagonal

// Move ...
type Move struct {
	state    uint32
	previous *Move
}

// GetPlayer ...
func (a *Move) GetPlayer() common.Player {
	if a.previous == nil {
		return common.NoPlayer{}
	}
	if isPlayer1(a.state) {
		return common.NewPlayer(0)
	}
	return common.NewPlayer(1)
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
	isNextPlayerP1 := !isPlayer1(a.state)
	board := getCombinedBoard(a.state)
	for i := 0; i < 9; i++ {
		bit := uint32(1) << i
		if board&bit == 0 {
			if !isNextPlayerP1 {
				bit = bit << 9
			}

			newState := flipPlayer(a.state) | bit
			m := Move{previous: a, state: newState}
			available = append(available, &m)
		}
	}

	return
}

// GetGameStatus ...
func (a *Move) GetGameStatus() common.GameStatus {
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
func (a *Move) BoardString() string {
	board1 := getBoard(a.state, true)
	board2 := getBoard(a.state, false)

	var b [9]string
	for i := 0; i < 9; i++ {
		bit := uint32(1) << i
		if board1&bit > 0 {
			b[i] = "X"
		} else if board2&bit > 0 {
			b[i] = "O"
		} else {
			b[i] = " "
		}
	}

	return fmt.Sprintf("%s|%s|%s\n-----\n%s|%s|%s\n-----\n%s|%s|%s\n",
		b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7], b[8])
}

// MoveString represents the current move
func (a *Move) MoveString() string {
	diff := (a.state - (*(a.previous)).state) & 262143
	if !isPlayer1(a.state) {
		diff = diff >> 9
	}
	playedIndex := bits.TrailingZeros32(diff)
	return strconv.Itoa(playedIndex)
}

// PlayMovesByIndex ...
func (a *Move) PlayMovesByIndex(moves *[]int) common.Move {
	return common.PlayMovesByIndex(a, moves)
}

// PlayMovesByString ...
func (a *Move) PlayMovesByString(moves string) common.Move {
	return common.PlayMovesByString(a, moves)
}

// NewGame start a new TicToe game by returning the opening move
func NewGame() common.Move {
	return &Move{262144, nil}
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
