package action

import (
	"b8boost/backend/internal/adapters/api/middleware"
	"b8boost/backend/internal/usecase"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin/binding"
)

type UpcomingEventsAction struct {
	uc usecase.UpcomingEventsUseCase
}

func NewUpcomingEventsAction(uc usecase.UpcomingEventsUseCase) UpcomingEventsAction {
	return UpcomingEventsAction{uc: uc}
}

func (a UpcomingEventsAction) Execute(w http.ResponseWriter, r *http.Request) {
	var input usecase.UpcomingEventInput
	if err := binding.Default(r.Method, binding.MIMEPOSTForm).Bind(r, &input); err != nil {
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

		input.UserID = &userID
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
