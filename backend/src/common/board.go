package common

// Location a square on a board...
type Location byte

// Offboard is a location off the board.  Use whan adding or removing a piece
var Offboard Location = Location(255)

func (l Location) index() byte {
	return byte(l)
}

func (l *Location) String() string {
	if l == nil {
		return "."
	}

	if byte(*l) < 26 {
		return string(97 + byte(*l)) // lower case
	}

	if byte(*l) < 52 {
		return string(65 + byte(*l) - 26) // upper case
	}

	return string(48 + byte(*l) - 52) // numbers
}

// LocationFromString ...
func LocationFromString(s string) Location {
	var index byte
	value := s[0]
	if value > 96 {
		index = value - 97
	} else if value > 64 {
		index = value - 65 + 26
	} else {
		index = value - 48 + 52
	}

	return Location(index)
}

// Board a small square board
type Board byte

// At ...
func (b Board) At(x byte, y byte) Location {
	//
	//             X
	//     ---------------
	//     |  0  1  2  3  4
	//     |  5  6  7  8  9
	//   Y | 10 11 12 13 14
	//     | 15 16 17 18 19
	//     | 20 21 22 23 24
	//

	index := x + b.Width()*y
	return Location(index)
}

func (b Board) Width() byte {
	return byte(b)
}

func (b Board) Height() byte {
	return byte(b)
}

func (b Board) X(l Location) byte {
	return l.index() % b.Width()
}

func (b Board) Y(l Location) byte {
	return l.index() / b.Width()
}

func (b Board) IsEdge(l Location) bool {
	x := b.X(l)
	if x == 0 || x == b.Width()-1 {
		return true
	}

	y := b.Y(l)
	return y == 0 || y == b.Height()-1
}

func (b Board) Up(l Location) (OnBoard bool, NewLocation Location) {
	if b.Y(l) == 0 {
		return offBoard()
	} else {
		return onBoard(l.index() - b.Width())
	}
}

func (b Board) Down(l Location) (OnBoard bool, NewLocation Location) {
	if b.Y(l) == b.Height()-1 {
		return offBoard()
	} else {
		return onBoard(l.index() + b.Width())
	}
}

func (b Board) Left(l Location) (OnBoard bool, NewLocation Location) {
	if b.X(l) == 0 {
		return offBoard()
	} else {
		return onBoard(byte(l) - 1)
	}
}

func (b Board) Right(l Location) (OnBoard bool, NewLocation Location) {
	if b.X(l) == b.Width()-1 {
		return offBoard()
	} else {
		return onBoard(byte(l) + 1)
	}
}

func offBoard() (onBoard bool, newLocation Location) {
	onBoard = false
	newLocation = Location(0)
	return
}

func onBoard(index byte) (OnBoard bool, NewLocation Location) {
	OnBoard = true
	NewLocation = Location(index)
	return
}
