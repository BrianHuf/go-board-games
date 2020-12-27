package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Basic(t *testing.T) {
	b := Board(5)
	
	corner := b.At(0,0)
	assert.Equal(t, "a", corner.String())
	assert.Equal(t, true, b.IsEdge(corner))

	corner2 := LocationFromString(corner.String())
	assert.Equal(t, corner, corner2)

	t.Log("Passed")
}
