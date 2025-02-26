package action

import (
	"b8boost/backend/internal/usecase"
	"encoding/json"
	"net/http"
)

type UpcomingEventsAction struct {
	uc usecase.UpcomingEventsUseCase
}

func NewUpcomingEventsAction(uc usecase.UpcomingEventsUseCase) UpcomingEventsAction {
	return UpcomingEventsAction{uc: uc}
}

func (a UpcomingEventsAction) Execute(w http.ResponseWriter, r *http.Request) {
	var input usecase.UpcomingEventInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	output, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, output)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
