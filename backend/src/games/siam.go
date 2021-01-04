package games

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

var board = common.Board(5)

// Direction ...
type Direction byte

const (
	directionUp Direction = iota
	directionDown
	directionLeft
	directionRight
	directionNone
)

func (d Direction) String() string {
	return [...]string{"U", "D", "L", "R"}[d]
}

// Opposite up -> down, left -> right, etc..
func (d Direction) Opposite() (opposite Direction) {
	if d == directionUp {
		opposite = directionDown
	} else if d == directionDown {
		opposite = directionUp
	} else if d == directionLeft {
		opposite = directionRight
	} else if d == directionRight {
		opposite = directionLeft
	} else if d == directionNone {
		opposite = directionNone
	} else {
		panic("invalid direction")
	}
	return
}

// Piece ...
type Piece byte

// String ...
func (p Piece) String() string {
	return [...]string{"(empty)", "Mountain",
		"P1 Up", "P1 Down", "P1 Left", "P1 Right",
		"P2 Up", "P2 Down", "P2 Left", "P2 Right"}[p]
}

// Char ...
func (p Piece) Char() string {
	if p == constantPieceEmpty {
		return "."
	}
	return [...]string{".", "M",
		"U", "D", "L", "R",
		"u", "d", "l", "r"}[p]
}

var (
	constantPieceEmpty    = Piece(0)
	constantPieceMountain = Piece(1)
)

// PiecesByPlayer ...
func PiecesByPlayer(p common.Player) []Piece {
	if p == common.Player1 {
		return []Piece{
			Piece(2), Piece(3), Piece(4), Piece(5)}
	} else if p == common.Player2 {
		return []Piece{
			Piece(6), Piece(7), Piece(8), Piece(9)}
	}

	return nil
}

// FromPlayer ...
func (p Piece) FromPlayer(player common.Player) bool {
	return (player == common.Player1 && byte(p) > 1 && byte(p) < 6) ||
		(player == common.Player2 && byte(p) > 5)
}

// GetDirection return Up/Down/Left/Right
func (p Piece) GetDirection() Direction {
	if p == 2 || p == 6 {
		return directionUp
	} else if p == 3 || p == 7 {
		return directionDown
	} else if p == 4 || p == 8 {
		return directionLeft
	} else if p == 5 || p == 9 {
		return directionRight
	}
	return directionNone
}

// BoardState 5x5 Board
type BoardState [25]Piece

// Move ...
type Move struct {
	playedBy common.Player
	previous *Move
	board    BoardState
	fromDir  Direction
	fromLoc  common.Location
	toLoc    common.Location
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
func (a *Move) GetJSON() interface{} {
	status := a.GetGameStatus()

	var winner string = ""
	var nextPlayer string = ""
	if status.IsDone() {
		if status.GetWinner() != common.PlayerNoOne {
			winner = status.GetWinner().String()
		}
	} else {
		nextPlayer = a.GetPlayer().Next().String()
	}

	return StateJSON{
		State:      a.BoardString(),
		LastMove:   a.MoveString(),
		IsDone:     status.IsDone(),
		Winner:     winner,
		NextPlayer: nextPlayer}
}

// GetPlayer ...
func (a *Move) GetPlayer() common.Player {
	return a.playedBy
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
	nextPlayer := a.GetPlayer().Next()
	freePieces := 5

	var addPieceMoves []common.Move
	var mountainCount int = 0

	for index, piece := range a.board {
		location := common.Location(byte(index))
		x := board.X(location)
		y := board.Y(location)

		if piece == constantPieceMountain {
			mountainCount++
		}

		// ADD PIECES
		if x == 0 {
			addPieceMoves = a.tryAdd(addPieceMoves, location, 3)
		} else if x == board.Width()-1 {
			addPieceMoves = a.tryAdd(addPieceMoves, location, 2)
		}

		if (x != 0 && x != board.Width()-1) || piece != constantPieceEmpty { // avoid duplicate corner moves
			if y == 0 {
				addPieceMoves = a.tryAdd(addPieceMoves, location, 1)
			} else if y == board.Height()-1 {
				addPieceMoves = a.tryAdd(addPieceMoves, location, 0)
			}
		}

		if piece.FromPlayer(nextPlayer) && board.IsEdge(location) {
			// REMOVE piece
			newBoard := a.board
			newBoard[index] = constantPieceEmpty

			newMove := Move{
				playedBy: nextPlayer,
				previous: a,
				board:    newBoard,
				fromLoc:  location,
				toLoc:    common.Offboard}

			available = append(available, &newMove)
		}

		if piece != constantPieceEmpty && piece.FromPlayer(nextPlayer) {
			freePieces--
			// spin
			for _, newPiece := range PiecesByPlayer(nextPlayer) {
				if newPiece == piece {
					continue
				}
				var newBoard = a.board
				newBoard[index] = newPiece

				var newMove = Move{
					playedBy: nextPlayer,
					previous: a,
					board:    newBoard,
					fromLoc:  location,
					toLoc:    location}

				available = append(available, &newMove)
			}

			// push/move
			available = a.tryMovePush(available, location, 0)
			available = a.tryMovePush(available, location, 1)
			available = a.tryMovePush(available, location, 2)
			available = a.tryMovePush(available, location, 3)
		}
	}

	if mountainCount < 3 {
		available = []common.Move{}
		return
	}

	if freePieces > 0 {
		available = append(available, addPieceMoves...)
	}

	return
}

func (a *Move) tryAdd(available []common.Move, to common.Location, dir Direction) []common.Move {
	validMove := a.canMovePush(common.Offboard, to, dir)
	if validMove {
		templateBoard := a.applyMovePush(common.Offboard, to, dir)

		nextPlayer := a.GetPlayer().Next()
		for _, newPiece := range PiecesByPlayer(nextPlayer) {
			var newBoard = templateBoard
			newBoard[to] = newPiece

			var newMove = Move{
				playedBy: nextPlayer,
				previous: a,
				board:    newBoard,
				fromDir:  dir,
				fromLoc:  common.Offboard,
				toLoc:    to}

			available = append(available, &newMove)
		}
	}

	return available
}

func (a *Move) tryMovePush(available []common.Move, from common.Location, dir Direction) []common.Move {
	onBoard, to := move(from, dir)
	if onBoard {
		validMove := a.canMovePush(from, to, dir)
		if validMove {
			templateBoard := a.applyMovePush(from, to, dir)

			nextPlayer := a.GetPlayer().Next()
			for _, newPiece := range PiecesByPlayer(nextPlayer) {
				var newBoard = templateBoard
				newBoard[to] = newPiece

				var newMove = Move{
					playedBy: nextPlayer,
					previous: a,
					board:    newBoard,
					fromLoc:  from,
					toLoc:    to}

				available = append(available, &newMove)
			}
		}
	}
	return available
}

func (a *Move) canMovePush(from common.Location, to common.Location, dir Direction) (canMove bool) {
	var onBoard bool

	power := byte(0)
	rocks := byte(0)

	current := from

	for {
		if current != common.Offboard && a.board[current] == constantPieceEmpty {
			canMove = true
			return
		}

		if current == common.Offboard || a.board[current].GetDirection() == dir {
			power++
		} else if a.board[current].GetDirection().Opposite() == dir {
			if power < 2 {
				canMove = false
				return
			}
			power--
		} else if a.board[current] == constantPieceMountain {
			rocks++
			if rocks > power {
				canMove = false
				return
			}
		} else if power == 0 && current == to && a.board[current] != constantPieceEmpty {
			canMove = false
			return
		}

		if current == common.Offboard {
			current = to
		} else {
			onBoard, current = move(current, dir)

			if !onBoard {
				canMove = true
				return
			}
		}
	}
}

func (a *Move) applyMovePush(from common.Location, to common.Location, dir Direction) BoardState {
	replaced := constantPieceEmpty

	newBoard := a.board

	if from == common.Offboard {
		replaced = PiecesByPlayer(a.GetPlayer().Next())[0]
	} else {
		replaced = newBoard[from]
		newBoard[from] = constantPieceEmpty
	}

	current := from
	for {
		var onBoard bool
		var next common.Location
		if current == common.Offboard {
			onBoard = true
			next = to
		} else {
			onBoard, next = move(current, dir)
		}

		if !onBoard {
			return newBoard
		}

		last := newBoard[next]
		newBoard[next] = replaced

		if last == constantPieceEmpty {
			return newBoard
		}
		replaced = last
		current = next
	}
}

func move(location common.Location, dir Direction) (onBoard bool, newLocation common.Location) {
	var x = board.X(location)
	var y = board.Y(location)

	if dir == 0 {
		if y == 0 {
			return false, location
		}
		return true, board.At(x, y-1)
	} else if dir == 1 {
		if y == board.Height()-1 {
			return false, location
		}
		return true, board.At(x, y+1)
	} else if dir == 2 {
		if x == 0 {
			return false, location
		}
		return true, board.At(x-1, y)
	} else if dir == 3 {
		if x == board.Width()-1 {
			return false, location
		}
		return true, board.At(x+1, y)
	}
	return false, location
}

func (a *Move) addMove(available []common.Move, from common.Location, to common.Location, dir Direction) []common.Move {
	nextPlayer := a.GetPlayer().Next()
	for _, newPiece := range PiecesByPlayer(nextPlayer) {
		var newBoard = a.board
		newBoard[from] = constantPieceEmpty
		newBoard[to] = newPiece

		var newMove = Move{
			playedBy: nextPlayer,
			previous: a,
			board:    newBoard,
			fromDir:  dir,
			fromLoc:  common.Offboard,
			toLoc:    to}

		available = append(available, &newMove)
	}

	return available
}

// GetGameStatus ...
func (a *Move) GetGameStatus() common.GameStatus {
	var mountainCount int = 0
	for _, piece := range a.board {
		if piece == constantPieceMountain {
			mountainCount++
		}
	}

	if mountainCount == 3 {
		return common.NewGameStatusInProgress()
	}

	return common.NewGameStatusWinner(a.playedBy)
}

// BoardString CLI based represention of game give current move
func (a *Move) BoardString() string {
	var b bytes.Buffer
	for _, p := range a.board {
		if p == constantPieceEmpty {
			b.WriteString(constantPieceEmpty.Char())
		} else {
			b.WriteString(p.Char())
		}
	}

	return b.String()
}

// MoveString represents the current move
func (a *Move) MoveString() string {
	if a.fromLoc == common.Offboard {
		return a.fromDir.String() + a.toLoc.String() + a.toPiece().Char()
	}
	return a.fromLoc.String() + a.toLoc.String() + a.toPiece().Char()
}

func (a *Move) toPiece() Piece {
	if a.toLoc == common.Offboard {
		return constantPieceEmpty
	}

	index := int(a.toLoc)
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
	initialBoard := BoardState{}
	initialBoard[11] = constantPieceMountain
	initialBoard[12] = constantPieceMountain
	initialBoard[13] = constantPieceMountain

	return &Move{
		playedBy: common.PlayerNoOne,
		previous: nil,
		board:    initialBoard,
		fromLoc:  common.Offboard,
		toLoc:    common.Offboard}
}
