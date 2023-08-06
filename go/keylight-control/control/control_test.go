package control

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	assert.Equal(t, Echo("Test"), "Testing")
}
