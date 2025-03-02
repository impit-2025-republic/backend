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
			DocID:   1,
			Title:   "Деятельность фонда",
			Content: "Фонд развития инноваций Республики Саха (Якутия) занимается поддержкой инновационных проектов и стартапов на ранних стадиях развития. Его основная деятельность включает:Финансирование стартапов – фонд предоставляет гранты, инвестиции и другие формы поддержки для предпринимателей, развивающих инновационные технологии.Акселерационные программы – фонд организует акселераторы, направленные на развитие стартапов, обучение предпринимателей и помощь в выходе на рынок. Например, программа B8 длится 3 месяца и помогает стартапам ускоренно развиваться.Поддержка импортозамещения – фонд финансирует проекты, связанные с разработкой отечественного программного обеспечения, в том числе для алмазодобывающей отрасли.Продвижение передовых технологий – в 2024 году фонд запустил программу по внедрению искусственного интеллекта в IT-продукты.Нетворкинг и менторство – фонд помогает стартапам находить партнеров, инвесторов и консультантов.Развитие инновационной экосистемы – фонд работает над созданием условий для появления и роста технологических компаний в регионе.Таким образом, фонд действует как ключевой инструмент поддержки технологического предпринимательства в Якутии.",
		},
		{
			DocID:   2,
			Title:   "Организационная структура",
			Content: "Фонд развития инноваций Республики Саха (Якутия) имеет следующую организационную структуру:Руководство:Директор: Птицына Вера Петровна, назначена на должность 9 ноября 2023 года. RUSPROFILE.RU Учредители:АО «Венчурная компания „Якутия» и ГАУ «Технопарк „Якутия»: эти организации совместно учредили фонд для поддержки стартапов на предпосевной стадии. INNOVATIONFUND14.RU Дочерние организации: Фонд имеет в своей структуре несколько дочерних компаний, каждая из которых выполняет определенные функции:ООО «СОФТВЭЙ++»: специализируется на разработке программного обеспечения. ООО «БЛ»: занимается предоставлением бизнес-услуг. ООО «СЕВЕРНЫЕ СКАЗКИ»: фокусируется на проектах в сфере культуры и искусства. Каждая из этих дочерних компаний способствует достижению целей фонда, направленных на развитие инновационной экосистемы в регионе.К сожалению, более подробная информация о внутренней структуре и подразделениях фонда в доступных источниках не представлена.",
		},
		{
			DocID:   3,
			Title:   "История фонда",
			Content: "Фонд развития инноваций Республики Саха (Якутия) был создан 19 декабря 2018 года совместными усилиями АО «Венчурная компания „Якутия» и ГАУ «Технопарк „Якутия». Его основная цель — поддержка стартапов на предпосевной стадии, особенно в сфере импортозамещения программного обеспечения. Фонд построен по модели эндаумент-фонда, что позволяет ему предоставлять финансовую и инфраструктурную поддержку, а также менторство и возможности для нетворкинга местным стартапам. INNOVATIONFUND14.RUЗа первые пять лет своей деятельности фонд стал ядром инновационной экосистемы Якутии, реализуя проекты, направленные на поддержку технологических предпринимателей. Среди таких проектов — акселерационная программа «Б8» и «Арктическая стартап-экспедиция: Дальний Восток и Арктика России». TPYKT.RUВ январе 2023 года АК «Алроса» заключила соглашение с фондом о пожертвовании 100 млн рублей в течение трех лет для поддержки проектов, связанных с импортозамещением программного обеспечения, особенно в алмазодобывающей отрасли. TADVISER.RUВ марте 2024 года фонд запустил акселератор стартапов, ориентированный на внедрение искусственного интеллекта в ИТ-продукты, что подчеркивает его стремление к развитию передовых технологий в регионе. CNEWS.RUТаким образом, Фонд развития инноваций Республики Саха (Якутия) играет ключевую роль в поддержке и развитии инновационной деятельности в регионе, способствуя созданию благоприятной среды для технологических стартапов и предпринимателей.",
		},
		{
			DocID:   4,
			Title:   "Директор фонда",
			Content: "Вера Петровна Птицына занимает должность директора Фонда развития инноваций Республики Саха (Якутия) с 9 ноября 2023 года. До этого она работала заместителем генерального директора АО «Венчурная компания „Якутия'». В рамках своей профессиональной деятельности Вера Петровна активно участвует в мероприятиях, направленных на развитие инновационной экосистемы региона. Например, 4 сентября 2024 года она выступала на сессии «Уйти в IT, или как построить цифровой суверенитет на базе кадрового потенциала» в рамках международного туристского форума TRAVEL HUB. Путешествуй! ROSCONGRESS.ORGПод руководством Веры Петровны фонд продолжает реализовывать программы, направленные на поддержку стартапов и развитие инноваций в Республике Саха (Якутия).",
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

	res, err := uc.ai.MakeVLLMRequest(messages, 0.3)
	if err != nil {
		return LLMChatOutput{}, err
	}
	return LLMChatOutput{
		Response: res,
	}, nil
}
