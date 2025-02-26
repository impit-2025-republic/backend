package service

import (
	"b8boost/backend/internal/entities"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type ErpSmartService struct {
	url         string
	accessToken string
	eventRepo   entities.EventRepo
}

func NewErpSmartService(
	eventRepo entities.EventRepo,
	accesstoken string,
	url string,
) ErpSmartService {
	return ErpSmartService{
		eventRepo:   eventRepo,
		url:         url,
		accessToken: accesstoken,
	}
}

type ErpTask struct {
	ID          int        `json:"id"`
	ProjectID   int        `json:"projectId"`
	BoardListID int        `json:"boardListId"`
	Name        string     `json:"name"`
	Priority    int        `json:"priority"`
	Description *string    `json:"description"`
	StartDate   time.Time  `json:"startDate"`
	DueDate     *time.Time `json:"dueDate"`
	Private     bool       `json:"private"`
	OwnerID     int        `json:"ownerId"`
	// UserID           *int        `json:"userId"`
	Closed           bool       `json:"closed"`
	ClosedAt         *time.Time `json:"closedAt"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	CommentsCount    string     `json:"commentsCount"`
	AttachmentsCount string     `json:"attachmentsCount"`
	Fio              string     `json:"fio"`
	Username         *string    `json:"username"`
	LabelNames       string     `json:"labelNames"`
	LabelIDs         string     `json:"labelIds"`
	MembersIDs       string     `json:"membersIds"`
	BoardListName    string     `json:"boardListName"`
	ProjectName      string     `json:"projectName"`
	LastComment      *string    `json:"lastComment"`
}

func (s ErpSmartService) Sync() {

	client := &http.Client{}

	req, err := http.NewRequest("GET", s.url, nil)
	if err != nil {
		log.Fatalf("Ошибка при создании запроса: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	cookie := &http.Cookie{
		Name:  "access_token",
		Value: s.accessToken,
		Path:  "/",
	}
	req.AddCookie(cookie)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Неожиданный код ответа: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Ошибка при чтении ответа: %v", err)
	}

	var tasks []ErpTask
	err = json.Unmarshal(body, &tasks)
	if err != nil {
		log.Fatalf("Ошибка при разборе JSON: %v", err)
	}

	var events []entities.Event
	for _, task := range tasks {
		status := "open"
		if task.Closed {
			status := "closed"
			_ = status
		}

		var desc string
		if task.Description != nil {
			desc = *task.Description
		}

		eventType := "task"

		event := entities.Event{
			EventID:     int64(task.ID),
			EventName:   task.Name,
			CreatedAt:   task.CreatedAt,
			Description: desc,
			Title:       task.Name,
			StartDs:     &task.StartDate,
			EndDs:       task.DueDate,
			Status:      &status,
			EventType:   &eventType,
			Coin:        float64(task.Priority),
			ErpID:       &task.ID,
		}

		events = append(events, event)
	}
	for _, event := range events {
		if event.ErpID != nil {
			existingEvent, err := s.eventRepo.FindByErpID(*event.ErpID)
			if err != nil && err == gorm.ErrRecordNotFound {
				originalID := existingEvent.EventID

				existingEvent.EventID = originalID

				err = s.eventRepo.Update(existingEvent)
				if err != nil {
					log.Printf("Ошибка при обновлении события с ErpID %d: %v", *event.ErpID, err)
					continue
				}
				log.Printf("Обновлено событие с ErpID %d", *event.ErpID)
			} else {
				err = s.eventRepo.Create(event)
				if err != nil {
					log.Printf("Ошибка при создании нового события с ErpID %d: %v", *event.ErpID, err)
					continue
				}
				log.Printf("Создано новое событие с ErpID %d", *event.ErpID)
			}
		}
	}
}
