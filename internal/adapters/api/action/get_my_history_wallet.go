package action

import (
	"b8boost/backend/internal/adapters/api/middleware"
	"b8boost/backend/internal/usecase"
	"net/http"
	"strconv"
)

type GetMyHistoryWalletAction struct {
	uc usecase.GetMyHistoryWalletUseCase
}

func NewGetMyHistoryWalletAction(uc usecase.GetMyHistoryWalletUseCase) GetMyHistoryWalletAction {
	return GetMyHistoryWalletAction{uc: uc}
}

func (a GetMyHistoryWalletAction) Execute(w http.ResponseWriter, r *http.Request) {
	var input usecase.GetMyHistoryWalletInput
	ctx := r.Context()
	userIdStr := ctx.Value(middleware.UserIDKey).(string)

	userID, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	input.UserID = userID

	output, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, output)
}
