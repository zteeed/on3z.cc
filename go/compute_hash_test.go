package swagger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveCharacters(t *testing.T) {
	t.Run("SuccessCharacterNotFound", func(t *testing.T) {
		input := "bonjour"
		res := removeCharacters(input, "a")
		assert.Equal(t, input, res)
	})
	t.Run("SuccessCharacterReplaced", func(t *testing.T) {
		input := "bonjour"
		res := removeCharacters(input, "b")
		assert.Equal(t, "onjour", res)
	})
	t.Run("SuccessCharacterReplacedTwice", func(t *testing.T) {
		input := "bonjour"
		res := removeCharacters(input, "o")
		assert.Equal(t, "bnjur", res)
	})
	t.Run("SuccessMultipleCharactersReplaced", func(t *testing.T) {
		input := "bonjour"
		res := removeCharacters(input, "bo")
		assert.Equal(t, "njur", res)
	})
}
