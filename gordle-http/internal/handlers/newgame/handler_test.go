package newgame

import (
	"gordle-http/internal/core"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type gameAdderStub struct {
	err error
}

func (g gameAdderStub) Add(_ core.Game) error {
	return g.err
}

func TestHandle(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/games", nil)

	require.NoError(t, err)

	recorder := httptest.NewRecorder()

	handlerFunc := Handler(gameAdderStub{})

	handlerFunc(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"id": "","attempts_left": 0,"guesses": null,"word_length": 0,"status": ""}`, recorder.Body.String())
}
