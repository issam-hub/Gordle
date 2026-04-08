package guess

import (
	"encoding/json"
	"gordle-http/internal/api"
	"net/http"
)

func Handle(w http.ResponseWriter, req *http.Request) {
	r := api.GuessRequest{}

	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
