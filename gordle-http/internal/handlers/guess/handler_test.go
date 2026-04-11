package guess

import (
	"gordle-http/internal/api"
	"gordle-http/internal/core"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type gameGuessStub struct {
	game core.Game
	err  error
}

func (g gameGuessStub) Modify(id string, r api.GuessRequest) (core.Game, error) {
	g.game = core.Game{
		ID: core.GameID(id),
	}
	return g.game, g.err
}

func TestHandle(t *testing.T) {
	body := strings.NewReader(`{"guess":"thing"}`)
	req, err := http.NewRequest(http.MethodPut, "/games", body)

	require.NoError(t, err)

	req.SetPathValue("id", "1")

	recorder := httptest.NewRecorder()

	handlerFunc := Handler(gameGuessStub{})

	handlerFunc(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, `{"id": "1","attempts_left": 0,"guesses": null,"word_length": 0,"status": ""}`, recorder.Body.String())
}
