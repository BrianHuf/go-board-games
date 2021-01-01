package common

import "fmt"

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
func PlayMovesByString(startingMove Move, moves string, chunkSize int) (nextMove Move) {
	if moves == "-" {
		return startingMove
	}

	nextMove = startingMove

	if chunkSize == 1 {
		for _, move := range moves {
			lookFor := string(move)
			lookInMoves := nextMove.NextAvailableMoves()
			nextMove = findMove(lookFor, &lookInMoves)
		}
	} else {
		for _, move := range chunks(moves, chunkSize) {
			lookFor := string(move)
			lookInMoves := nextMove.NextAvailableMoves()
			nextMove = findMove(lookFor, &lookInMoves)
		}
	}

	return
}

// https://stackoverflow.com/a/61469854/5468964
func chunks(s string, chunkSize int) []string {
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	chunk := make([]rune, chunkSize)
	len := 0
	for _, r := range s {
		chunk[len] = r
		len++
		if len == chunkSize {
			chunks = append(chunks, string(chunk))
			len = 0
		}
	}
	if len > 0 {
		chunks = append(chunks, string(chunk[:len]))
	}
	return chunks
}

func findMove(lookFor string, moves *[]Move) Move {
	for _, move := range *moves {
		check := move.MoveString()
		if check == lookFor {
			return move
		}
	}
	panic(fmt.Sprintf("unable to find move %s", lookFor))
}
