package tictactoe

import (
	"fmt"
	"math/bits"
	"strconv"

	"me.dev/go-board-game/common"
)

// Move ...
type Move struct {
	state    uint32
	previous *Move
}

// GetPlayer ...
func (a Move) GetPlayer() common.Player {
	if a.previous == nil {
		return common.GamePlayer{}
	}
	return common.TwoPlayer{IsPlayer1: isPlayer1(a.state)}
}

// GetPrevious ...
func (a Move) GetPrevious() common.Move {
	if a.previous == nil {
		return nil
	}
	return *(a.previous)
}

// String ...
func (a Move) String() string {
	return "board"
}

// NextAvailableMoves ...
func (a Move) NextAvailableMoves() []common.Move {
	var available []common.Move

	p1 := isPlayer1(a.state)
	board := getCombinedBoard(a.state)
	for i := 0; i < 9; i++ {
		bit := uint32(1) << i
		if board&bit == 0 {
			if p1 {
				bit = bit << 9
			}

			newState := a.state | bit
			if a.state != 0 {
				newState = flipPlayer(newState)
			}
			m := Move{previous: &a, state: newState}
			available = append(available, m)
		}
	}

	return available
}

// GetWinner ...
func (a Move) GetWinner() common.Player {
	if a.previous == nil {
		return nil
	}

	p1 := isPlayer1(a.state)
	if isWinningBoard(getBoard(a.state, p1)) {
		return common.TwoPlayer{IsPlayer1: p1}
	}

	return nil
}

// BoardString CLI based represention of game give current move
func (a Move) BoardString() string {
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
func (a Move) MoveString() string {
	diff := (a.state - (*(a.previous)).state) & 262143
	if !isPlayer1(a.state) {
		diff = diff >> 9
	}
	playedIndex := bits.TrailingZeros32(diff)
	return strconv.Itoa(playedIndex)
}

// PlayMovesByIndex ...
func (a Move) PlayMovesByIndex(moves []int) common.Move {
	return common.PlayMovesByIndex(a, moves)
}

// PlayMovesByString ...
func (a Move) PlayMovesByString(moves string) common.Move {
	return common.PlayMovesByString(a, moves)
}

// NewGame start a new TicToe game by returning the opening move
func NewGame() common.Move {
	return Move{0, nil}
}

// local functions
func isPlayer1(state uint32) bool {
	return state != 0 && (state&262144) == 0
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
	return board == 7 || board == 56 || board == 448 || // horizontal
		board == 73 || board == 584 || board == 4672 || // vertical
		board == 84 || board == 273 // diagonal
}
