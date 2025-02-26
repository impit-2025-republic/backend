package action

import (
	"b8boost/backend/internal/adapters/api/middleware"
	"b8boost/backend/internal/usecase"
	"net/http"
	"strconv"
)

type UserMeAction struct {
	uc usecase.UserMeUseCase
}

func NewUserMeAction(uc usecase.UserMeUseCase) UserMeAction {
	return UserMeAction{uc: uc}
}

func (a UserMeAction) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIdStr := ctx.Value(middleware.UserIDKey).(string)

	userID, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	output, err := a.uc.Execute(r.Context(), usecase.UserMeInput{
		UserID: userID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, output)
}
