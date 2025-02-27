package action

import (
	"b8boost/backend/internal/usecase"
	"net/http"
)

type FindProductAction struct {
	uc usecase.FindProductUseCase
}

func NewFindProductAction(uc usecase.FindProductUseCase) FindProductAction {
	return FindProductAction{uc: uc}
}

func (a FindProductAction) Execute(w http.ResponseWriter, r *http.Request) {
	output, err := a.uc.Execute(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, output)
}
