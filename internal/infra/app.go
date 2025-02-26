package infra

import (
	"b8boost/backend/config"
	"b8boost/backend/internal/infra/ai"
	"b8boost/backend/internal/infra/cron"
	"b8boost/backend/internal/infra/database"
	"b8boost/backend/internal/infra/jwt"
	"b8boost/backend/internal/infra/ldap"
	"b8boost/backend/internal/infra/router"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"gorm.io/gorm"
)

type app struct {
	ai     ai.Vllm
	cfg    config.Config
	db     *gorm.DB
	jwt    jwt.JWKSHandler
	router router.RouterHTTP
	ldap   ldap.LDAP
}

func parsePrivateKey(pemString string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	var ok bool
	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}

	return privateKey, nil
}

func Config(cfg config.Config) *app {
	return &app{cfg: cfg}
}

func (a *app) JWT() *app {
	privateKey, err := parsePrivateKey(a.cfg.JWTPrivateKey)
	if err != nil {
		panic(err)
	}
	a.jwt = jwt.NewJWKSHandler(privateKey)
	return a
}

func (a *app) Database() *app {
	a.db = database.NewGormDB(a.cfg)
	return a
}

func (a *app) Ldap() *app {
	a.ldap = ldap.NewLDAP(a.cfg)
	return a
}

func (a *app) Serve() *app {
	a.router = router.NewRouterHTTP(a.jwt, a.cfg.BotToken, a.ldap, a.db, a.ai)
	return a
}

func (a *app) Cron() *app {
	cron := cron.NewCron(a.db, a.ldap, a.cfg)
	cron.Start()
	return a
}

func (a *app) LLM() *app {
	a.ai = ai.NewVllm(a.cfg.VLLMUrl, a.cfg.VLLMApiKey)
	return a
}

func (a *app) Start() {
	a.router.Listen()
}
