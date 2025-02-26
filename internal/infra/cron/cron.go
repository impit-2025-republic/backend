package cron

import (
	"b8boost/backend/internal/adapters/repo"
	"b8boost/backend/internal/adapters/service"
	"b8boost/backend/internal/infra/ldap"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Cron struct {
	db   *gorm.DB
	ldap ldap.LDAP
}

func NewCron(db *gorm.DB, ldap ldap.LDAP) Cron {
	return Cron{
		db:   db,
		ldap: ldap,
	}
}

func (c Cron) Start() {
	checkout := cron.New(cron.WithLocation(time.Local))
	checkout.AddFunc("*/15 * * * *", func() {
		service := service.NewLDAPService(c.ldap, repo.NewUserRepo(c.db))
		service.Sync()
	})

	checkout.Start()
}
