package entities

import (
	"time"
)

type (
	UserRepo interface {
		GetByLdapID(entryUUID string) (User, error)
		GetByID(id uint) (User, error)
		GetAll() ([]User, error)
		Create(user User) (User, error)
		Update(user User) error
	}

	User struct {
		UserID      uint       `gorm:"primaryKey;column:user_id;autoIncrement"`
		TelegramID  *int       `gorm:"column:telegram_id"`
		Surname     *string    `gorm:"column:surname;type:varchar(255)"`
		Name        *string    `gorm:"column:name;type:varchar(255)"`
		LastSurname *string    `gorm:"column:last_surname;type:varchar(255)"`
		BirthDate   *time.Time `gorm:"column:birth_date;type:date"`
		Role        *string    `gorm:"column:role;type:varchar(100)"`
		CompanyID   *int       `gorm:"column:company_id"`
		CreatedAt   time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
		UpdatedAt   *time.Time `gorm:"column:updated_at;type:timestamp"`
		Description *string    `gorm:"column:description;type:text"`
		Avatar      *string    `gorm:"column:avatar;type:varchar(255)"`
		LastLogin   *time.Time `gorm:"column:last_login;type:timestamp"`
		IsOnline    *bool      `gorm:"column:is_online"`
		Email       *string    `gorm:"column:email;type:varchar(255)"`
		Phone       *string    `gorm:"column:phone;type:varchar(100)"`
		LdapID      string     `gorm:"column:ldap_id"`
	}
)

func (User) TableName() string {
	return "users"
}
