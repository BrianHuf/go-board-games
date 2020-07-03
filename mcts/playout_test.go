package mcts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"me.dev/go-board-game/games/tictactoe"
)

func Test_Basic(t *testing.T) {
	count := 1000
	var scores [9]float32
	for startingMove := 0; startingMove < 9; startingMove++ {
		move := tictactoe.NewGame().PlayMovesByIndex([]int{startingMove})

		var totalScore float32 = 0.0
		for i := 0; i < count; i++ {
			_, score := BasicRandomPlayout(move)
			score /= float32(count)
			totalScore += score
		}
		scores[startingMove] = totalScore
		t.Logf("score at %d -> %+.3f", startingMove, totalScore)
	}

	corners := (scores[0] + scores[2] + scores[6] + scores[8]) / 4.0
	centers := (scores[1] + scores[3] + scores[5] + scores[7]) / 4.0

	assert.InEpsilon(t, 0.130, corners, 0.01)
	assert.InEpsilon(t, -0.037, centers, 0.01)
	assert.InEpsilon(t, 0.285, scores[4], 0.01)

}
