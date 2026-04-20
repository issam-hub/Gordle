package status

import (
	"gordle-http/internal/core"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type gameGetterStub struct {
	game core.Game
	err  error
}

func (g gameGetterStub) Get(id string) (core.Game, error) {
	g.game = core.Game{
		ID: core.GameID(id),
	}
	return g.game, g.err
}

func TestHandle(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/games/", nil)
	require.NoError(t, err)

	req.SetPathValue("id", "1")

	recorder := httptest.NewRecorder()

	handlerFunc := Handler(gameGetterStub{})

	handlerFunc(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, `{"id": "1","attempts_left": 0,"guesses": null,"word_length": 5,"status": ""}`, recorder.Body.String())
}
