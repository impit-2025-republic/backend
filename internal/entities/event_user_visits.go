package entities

type (
	EventUserVisitRepo interface {
		Create(event EventUserVisit) error
		GetByEventIDAndUserID(eventID, userID uint) (EventUserVisit, error)
	}

	EventUserVisit struct {
		EventID           int    `gorm:"column:event_id"`
		UserID            int    `gorm:"column:user_id"`
		Visit             string `gorm:"column:visit;type:varchar(255)"`
		AchievementTypeID int    `gorm:"column:achievement_type_id"`
	}
)

func (EventUserVisit) TableName() string {
	return "event_user_visits"
}
