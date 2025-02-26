package action

import (
	"b8boost/backend/internal/usecase"
	"encoding/json"
	"io"
	"net/http"
)

type LLMChatAction struct {
	uc usecase.LLMChatUseCase
}

func NewLLMChatAction(uc usecase.LLMChatUseCase) LLMChatAction {
	return LLMChatAction{uc: uc}
}

func (a LLMChatAction) Execute(w http.ResponseWriter, r *http.Request) {
	var input usecase.LLMChatInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	output, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := output.Response

	if res == nil {
		http.Error(w, "Null response received", http.StatusInternalServerError)
		return
	}

	for key, values := range res.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(res.StatusCode)

	if res.Body != nil {
		defer res.Body.Close()
		io.Copy(w, res.Body)
	}
}
