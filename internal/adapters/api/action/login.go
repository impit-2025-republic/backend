package action

import (
	"b8boost/backend/internal/usecase"
	"net/http"
	"strings"
)

type LoginAction struct {
	uc usecase.LoginUsecase
}

func NewLoginAction(uc usecase.LoginUsecase) LoginAction {
	return LoginAction{uc: uc}
}

func (a LoginAction) Execute(w http.ResponseWriter, r *http.Request) {

	initData := strings.Split(r.Header.Get("Authorization"), " ")

	var input usecase.LoginInput

	if len(initData) == 2 {
		input.InitData = initData[1]
	}

	output, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output.Token))
}
