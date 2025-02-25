package action

import (
	"b8boost/backend/internal/usecase"
	"encoding/json"
	"net/http"
)

type FindJwtsAction struct {
	uc usecase.FindJwtsUseCase
}

func NewFindJwtsAction(uc usecase.FindJwtsUseCase) FindJwtsAction {
	return FindJwtsAction{
		uc: uc,
	}
}

func (a *FindJwtsAction) Execute(w http.ResponseWriter, r *http.Request) {
	jwts := a.uc.Execute(r.Context())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwts)
}
