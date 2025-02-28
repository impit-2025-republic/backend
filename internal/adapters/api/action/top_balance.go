package action

import (
	"b8boost/backend/internal/usecase"
	"net/http"
)

type TopBalanceAction struct {
	uc usecase.TopBalanceUseCase
}

func NewTopBalanceAction(uc usecase.TopBalanceUseCase) TopBalanceAction {
	return TopBalanceAction{uc: uc}
}

func (a TopBalanceAction) Execute(w http.ResponseWriter, r *http.Request) {

	output, err := a.uc.Execute(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, output)
}
