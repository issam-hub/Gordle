package status

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/games/", nil)
	require.NoError(t, err)

	req.SetPathValue("id", "1")

	recorder := httptest.NewRecorder()

	Handle(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, `{"id": "1","attempts_left": 0,"guesses": null,"word_length": 0,"status": ""}`, recorder.Body.String())
}
