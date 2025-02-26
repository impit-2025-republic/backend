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
			Title:   "Аутизм у детей, общая информация",
			Content: "Аутизм у детей – болезнь, связанная с нарушением развития определенных психических функций, которая проявляется различными трудностями в социальном взаимодействии ребенка с окружающим миром, навязчивыми двигательными привычками и другими состояниями. Чаще всего заболевание диагностируют у малышей до 3–4 лет, но первые признаки отклонений в некоторых случаях можно заметить уже на первом году жизни. Общего лечения патологии не существует: специалисты разрабатывают индивидуальные методы коррекции аутизма у детей, отталкиваясь от состояния конкретного пациента. Обратите внимание: в детском отделении «СМ-Клиника» лечение аутизма не проводится, осуществляется общая поддержка и консультативная помощь в рамках назначенной ранее коррекции. В последние десятилетия число малышей с расстройством аутического спектра (РАС) значительно увеличилось. Врачи, работающие с такими ребятами, уверены, что причиной резкого прироста выявленных случаев послужила смена критериев обследования, развитие диагностических методик и более тщательное изучение проблемы. Степень проявления симптоматики при аутизме у детей может существенно отличаться: от полной неспособности контактировать с другими людьми до определенных «странностей» в поведении, таких как навязчивые движения, слишком узкий круг интересов или нестандартная манера речи. Более склонны к расстройству аутического спектра мальчики: они сталкиваются с заболеванием в общей сложности в два-три раза чаще, чем девочки. Ученые объясняют это лучшими коммуникативными способностями женского пола, из-за чего слабо выраженные формы аутизма могут быть попросту не замечены. Оценкой симптомов, поиском возможных причин возникновения аутизма у детей и коррекцией признаков РАС занимаются в комплексе нейропсихологи, детские неврологи, логопеды, психиатры, а также педагоги-дефектологи и социальные службы, если требуется.",
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
