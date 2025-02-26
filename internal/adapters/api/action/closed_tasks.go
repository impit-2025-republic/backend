package action

import (
	"b8boost/backend/internal/usecase"
	"net/http"
)

type ClosedEventsAction struct {
	uc usecase.ClosedEventsUseCase
}

func NewClosedEventsAction(uc usecase.ClosedEventsUseCase) ClosedEventsAction {
	return ClosedEventsAction{uc: uc}
}

func (a ClosedEventsAction) Execute(w http.ResponseWriter, r *http.Request) {
	output, err := a.uc.Execute(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, output)
}
