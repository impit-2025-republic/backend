package action

import (
	"b8boost/backend/internal/usecase"
	"encoding/json"
	"net/http"
)

type CreateEventAction struct {
	uc usecase.CreateEventUseCase
}

func NewCreateEventAction(uc usecase.CreateEventUseCase) CreateEventAction {
	return CreateEventAction{uc: uc}
}

func (a CreateEventAction) Execute(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateEventInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(output)
}
