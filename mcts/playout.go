package mcts

import (
	"math/rand"

	"me.dev/go-board-game/common"
)

// PlayoutSelector function to select the next move
type PlayoutSelector func(currentMove common.Move) (nextMove common.Move)

// PlayoutScorer function to calculate a playout score (-1.0 bad - 1.0 good)
type PlayoutScorer func(finalMove common.Move, numMoves int, startingPlayer common.Player, winningPlayer common.Player) (score float32)

// BasicRandomPlayout randomly select up to 1000 moves and returns a score (-1.0 bad, 0.0 tie, 1.0 good)
func BasicRandomPlayout(m common.Move) (finalMove common.Move, score float32) {
	return BasicPlayout(m, 1000, RandomPlayoutSelector, BasicPlayoutScorer)
}

// BasicPlayout automatically select moves until game ends with winner or no more moves.  Return a score (-1.0 bad - 1.0 good)
func BasicPlayout(m common.Move, maxMoves int, selector PlayoutSelector, scorer PlayoutScorer) (finalMove common.Move, score float32) {
	startingPlayer := m.GetPlayer()
	currentMove := m

	for count := 1; count <= maxMoves; count++ {
		currentMove = selector(currentMove)
		status := currentMove.GetGameStatus()
		//fmt.Println(currentMove.BoardString() + status.String())
		if status.IsDone() {
			score = scorer(currentMove, count, startingPlayer, status.GetWinner())
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

// BasicPlayoutScorer simple score -1.0 if lost, 0.0 if tie, 1.0 if win
func BasicPlayoutScorer(finalMove common.Move, numMoves int, startingPlayer common.Player, winningPlayer common.Player) (score float32) {
	if winningPlayer == nil {
		score = 0.0
	} else if winningPlayer.Equals(startingPlayer) {
		score = 1.0
	} else {
		score = -1.0
	}
	return score
}
