package action

import (
	"b8boost/backend/internal/usecase"
	"encoding/json"
	"net/http"
)

type AdminVisitEventAction struct {
	uc usecase.AdminVisitEventUseCase
}

func NewAdminVisitEventAction(uc usecase.AdminVisitEventUseCase) AdminVisitEventAction {
	return AdminVisitEventAction{uc: uc}
}

func (a AdminVisitEventAction) Execute(w http.ResponseWriter, r *http.Request) {
	var input usecase.AdminVisitEventInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = a.uc.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
