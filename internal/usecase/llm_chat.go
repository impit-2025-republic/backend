package usecase

import (
	"b8boost/backend/internal/infra/ai"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	LLMChatUseCase interface {
		Execute(ctx context.Context, input LLMChatInput) (LLMChatOutput, error)
	}

	LLMChatInput struct {
		Promnt string `json:"promnt"`
	}

	LLMChatOutput struct {
		Response *http.Response
	}

	llmChatInteractor struct {
		ai ai.Vllm
	}
)

func NewLLmChatInteractor(ai ai.Vllm) LLMChatUseCase {
	return llmChatInteractor{
		ai: ai,
	}
}

func (uc llmChatInteractor) Execute(ctx context.Context, input LLMChatInput) (LLMChatOutput, error) {
	systemPrompt := "Your task is to answer the user's questions using only the information from the provided documents. Give two answers to each question: one with a list of relevant document identifiers and the second with the answer to the question itself, using documents with these identifiers."
	documents := []ai.Document{
		{
			DocID:   0,
			Title:   "Фонд развитии инновации - ФРИ",
			Content: "О фонде: Создан совместно АО «Венчурная компания «Якутия» и ГАУ «Технопарк «Якутия»» построен по модели эндаумент-фонда для оказания финансовой поддержки стартапам на предпосевной стадии. Цели: Способствовать развитию инновационной экосистемы Усилить инвестиционный потенциал местных стартапов через оказание финансовой поддержки, а также через проведение акселерационной программы Деятельность:  1. Б8 Акселератор - Длительность программы – 3 месяца 3. Инфраструктурная поддержка 4. Менторская поддержка 5. Нетворкинг 6. Финансовая поддержка В АО «Венчурная Компания “Якутия” открылась вакансия стажера, в «Фонд развития инноваций Республики Саха (Якутия) — менеджера по работе с ИТ-проектами.",
		},
		{
			DocID:   1,
			Title:   "Программа стажировок «РОСТ»",
			Content: "уникальная возможность для амбициозных начинающих специалистов начать свою карьеру в сфере инноваций и технологического предпринимательства! Мы уверены, что для того, чтобы наша компания была лучшей и постоянно росла, мы должны развивать людей, которые ведут ее к этим высотам.",
		},
	}
	docsJson, err := json.Marshal(documents)
	if err != nil {
		return LLMChatOutput{}, err
	}

	messages := []ai.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "documents", Content: string(docsJson)},
		{Role: "user", Content: input.Promnt},
	}

	indexes, err := uc.ai.MakeVLLMIndexes(messages, 0.0)
	if err != nil {
		return LLMChatOutput{}, err
	}
	fmt.Println("Indexes:", indexes)

	messages = append(messages, ai.Message{
		Role:    "assistant",
		Content: indexes,
	})

	res, err := uc.ai.MakeVLLMRequest(messages, 0.0)
	if err != nil {
		return LLMChatOutput{}, err
	}
	return LLMChatOutput{
		Response: res,
	}, nil
}
