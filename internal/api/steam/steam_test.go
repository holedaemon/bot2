package steam

import (
	"context"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestSteam(t *testing.T) {
	key := os.Getenv("STEAM_KEY")
	id := os.Getenv("STEAM_ID")
	assert.Assert(t, key != "" && id != "", "env variables not set")

	c, err := New(key)
	assert.NilError(t, err, "creating new client")

	ownedGames, err := c.GetOwnedGames(context.Background(), id)
	assert.NilError(t, err, "getting owned games")

	assert.Assert(t, len(ownedGames.Games) > 0, "owned games is blank")
}
