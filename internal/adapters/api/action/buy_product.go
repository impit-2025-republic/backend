package action

import (
	"b8boost/backend/internal/adapters/api/middleware"
	"b8boost/backend/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin/binding"
)

type BuyProductAction struct {
	uc usecase.BuyProductUseCase
}

func NewBuyProductAction(uc usecase.BuyProductUseCase) BuyProductAction {
	return BuyProductAction{uc: uc}
}

func (a BuyProductAction) Execute(w http.ResponseWriter, r *http.Request) {
	var input usecase.BuyProductInput
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

		input.UserID = uint(userID)
	}

	err := a.uc.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
