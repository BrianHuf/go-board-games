package mcts

import (
	"math/rand"

	"me.dev/go-board-game/common"
)

// BasicRandomPlayout randomly select up to 1000 moves and returns a score (0.0 bad, 0.5 tie, 1.0 good)
func BasicRandomPlayout(m common.Move) (finalMove common.Move, score float32) {
	return Playout(m, 1000, RandomPlayoutSelector, BasicPlayoutScorer)
}

// PlayoutSelector function to select the next move
type PlayoutSelector func(currentMove common.Move) (nextMove common.Move)

// PlayoutScorer function to calculate a playout score (0.0 bad - 1.0 good)
type PlayoutScorer func(finalMove common.Move, numMoves int, startingPlayer common.Player, winningPlayer common.Player) (score float32)

// Playout automatically select moves until game ends with winner or no more moves.  Return a score (0.0 bad - 1.0 good)
func Playout(m common.Move, maxMoves int, selector PlayoutSelector, scorer PlayoutScorer) (finalMove common.Move, score float32) {
	startingPlayer := m.GetPlayer()
	currentMove := m

	for count := 1; count <= maxMoves; count++ {
		nextMove := selector(currentMove)
		if nextMove == nil {
			score = scorer(currentMove, count-1, startingPlayer, nil)
			finalMove = currentMove
			return finalMove, score
		}

		currentMove = nextMove
		winningPlayer := currentMove.GetWinner()
		if winningPlayer != nil {
			score = scorer(currentMove, count, startingPlayer, winningPlayer)
			finalMove = currentMove
			return finalMove, score
		}
	}

	score = scorer(currentMove, maxMoves, startingPlayer, nil)
	finalMove = currentMove
	return finalMove, score
}

// RandomPlayoutSelector randomly select moves
func RandomPlayoutSelector(m common.Move) (nextMove common.Move) {
	moves := m.NextAvailableMoves()
	numMoves := len(moves)
	if numMoves == 0 {
		nextMove = nil
		return nextMove
	}

	index := rand.Intn(numMoves)
	nextMove = moves[index]
	return nextMove
}

// BasicPlayoutScorer simple score 0.0 if lost, 0.5 if tie, 1.0 if win
func BasicPlayoutScorer(finalMove common.Move, numMoves int, startingPlayer common.Player, winningPlayer common.Player) (score float32) {
	if winningPlayer == nil {
		score = 0.5
	} else if winningPlayer.Equals(startingPlayer) {
		score = 1.0
	} else {
		score = 0.0
	}
	return score
}
