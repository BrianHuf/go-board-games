package mcts

import (
	"math/rand"

	"me.dev/go-board-game/common"
)

// Playout automatically select moves until game ends with winner or no more moves.  Return a score (0.0 bad - 1.0 good)
func Playout(m common.Move, config Config) (score float32, finalMove common.Move) {
	currentMove := m
	startingPlayer := currentMove.GetPlayer()

	var status common.GameStatus = nil

	for count := 1; count <= config.PlayoutMaxMoves; count++ {
		status = currentMove.GetGameStatus()
		if status.IsDone() {
			score = config.PlayoutScorer(currentMove, count, startingPlayer, status.GetWinner())
			finalMove = currentMove
			return
		}

		if count < config.PlayoutMaxMoves {
			currentMove = config.PlayoutSelector(currentMove)
		}
	}

	score = config.PlayoutScorer(currentMove, config.PlayoutMaxMoves, startingPlayer, common.PlayerNoOne)
	finalMove = currentMove
	return
}

// BasicPlayoutSelector randomly select moves
func BasicPlayoutSelector(m common.Move) (nextMove common.Move) {
	moves := m.NextAvailableMoves()
	numMoves := len(moves)
	if numMoves == 0 {
		nextMove = nil
		return nextMove
	}

	index := rand.Intn(numMoves)
	nextMove = moves[index]
	return
}

// BasicPlayoutScorer simple score 0.0 if lost, 0.5 if tie, 1.0 if win
func BasicPlayoutScorer(finalMove common.Move, numMoves int, startingPlayer common.Player, winningPlayer common.Player) (score float32) {
	if winningPlayer == common.PlayerNoOne {
		score = 0.5
	} else if winningPlayer == startingPlayer {
		score = 1.0
	} else {
		score = 0.0
	}
	return score
}
