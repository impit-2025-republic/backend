package action

import (
	"b8boost/backend/internal/adapters/api/middleware"
	"b8boost/backend/internal/usecase"
	"encoding/json"
	"net/http"
	"strconv"
)

type CaseOpenAction struct {
	uc usecase.CaseOpenUseCase
}

func NewCaseOpenAction(uc usecase.CaseOpenUseCase) CaseOpenAction {
	return CaseOpenAction{uc: uc}
}

func (a CaseOpenAction) Execute(w http.ResponseWriter, r *http.Request) {
	var input usecase.CaseOpenInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx := r.Context()
	userIdStr, ok := ctx.Value(middleware.UserIDKey).(string)
	if ok {
		userID, err := strconv.Atoi(userIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		input.UserID = userID
	}

	output, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, output)
}
