package cron

import (
	"b8boost/backend/config"
	"b8boost/backend/internal/adapters/repo"
	"b8boost/backend/internal/adapters/service"
	"b8boost/backend/internal/infra/ldap"
	"b8boost/backend/internal/infra/tgbot"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Cron struct {
	db    *gorm.DB
	ldap  ldap.LDAP
	cfg   config.Config
	tgbot tgbot.TgBot
}

func NewCron(db *gorm.DB, ldap ldap.LDAP, cfg config.Config, tgbot tgbot.TgBot) Cron {
	return Cron{
		db:    db,
		ldap:  ldap,
		cfg:   cfg,
		tgbot: tgbot,
	}
}

func (c Cron) Start() {
	checkout := cron.New(cron.WithLocation(time.Local))
	checkout.AddFunc("*/1 * * * *", func() {
		service := service.NewLDAPService(
			c.ldap,
			repo.NewUserRepo(c.db),
			repo.NewUserWallet(c.db))
		service.Sync()
	})

	checkout.AddFunc("*/1 * * * *", func() {
		service := service.NewEventStatusService(
			repo.NewEventRepo(c.db),
			repo.NewEventUserVisits(c.db),
			repo.NewUserWallet(c.db),
			repo.NewUserWalletHistoryRepo(c.db),
			c.tgbot,
		)
		service.Start()
	})

	checkout.AddFunc("@every 30s", func() {
		service := service.NewErpSmartService(
			repo.NewEventRepo(c.db),
			c.cfg.ERPAccessToken,
			c.cfg.ERPURL,
		)
		service.Sync()
	})

	checkout.Start()
}
