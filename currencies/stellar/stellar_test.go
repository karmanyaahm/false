package stellar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStellar(t *testing.T) {

	t.Log("HI")
	assert.Equal(t, "HI", "HI")
	s := &Stellar{}
	s.Init("idk", func())
}
