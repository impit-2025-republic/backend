package action

import (
	"b8boost/backend/internal/adapters/api/middleware"
	"b8boost/backend/internal/usecase"
	"encoding/json"
	"net/http"
	"strconv"
)

type VisitEventAction struct {
	uc usecase.VisitEventUseCase
}

func NewVisitEventAction(uc usecase.VisitEventUseCase) VisitEventAction {
	return VisitEventAction{uc: uc}
}

func (a VisitEventAction) Execute(w http.ResponseWriter, r *http.Request) {
	var input usecase.VisitEventInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx := r.Context()
	userIdStr := ctx.Value(middleware.UserIDKey).(string)

	userID, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	input.UserID = userID

	err = a.uc.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
