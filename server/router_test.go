package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	assert.NotNil(t, New())
}
