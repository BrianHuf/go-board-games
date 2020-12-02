package common

// BasicMove provide some common helpers most games would benefit from ...
type BasicMove struct{}

// PlayMovesByIndex ...
func PlayMovesByIndex(startingMove Move, moves *[]int) Move {
	var ret Move = startingMove
	for _, moveIndex := range *moves {
		ret = ret.NextAvailableMoves()[moveIndex]
	}
	return ret
}

// PlayMovesByString ...
func PlayMovesByString(startingMove Move, moves string) (nextMove Move) {
	if moves == "-" {
		return startingMove
	}

	nextMove = startingMove
	for _, move := range moves {
		lookFor := string(move)
		lookInMoves := nextMove.NextAvailableMoves()
		nextMove = findMove(lookFor, &lookInMoves)
	}
	return
}

func findMove(lookFor string, moves *[]Move) Move {
	for _, move := range *moves {
		check := move.MoveString()
		if check == lookFor {
			return move
		}
	}
	return nil
}
