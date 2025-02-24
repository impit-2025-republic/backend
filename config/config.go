package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	JWTPrivateKey    string `env:"JWT_PRIVATE_KEY" env-required:"true"`
	DatabaseHost     string `env:"DATABASE_HOST" env-default:"localhost"`
	DatabasePort     int    `env:"DATABASE_PORT" env-default:"5432"`
	DatabaseUser     string `env:"DATABASE_USER" env-default:"postgres"`
	DatabasePassword string `env:"DATABASE_PASSWORD" env-default:"postgres"`
	DatabaseDB       string `env:"DATABASE_DB" env-default:"b8st"`
	BotToken         string `env:"BOT_TOKEN" env-required:"true"`

	LDAPServer     string `env:"LDAP_SERVER" env-default:"localhost"`
	LDAPPort       string `env:"LDAP_PORT" env-default:"389"`
	LDAPBindDN     string `env:"LDAP_BIND_DN" env-default:"cn=admin,dc=example,dc=com"`
	LDAPBindPass   string `env:"LDAP_BIND_PASS" env-default:"admin"`
	LDAPBaseDN     string `env:"LDAP_BASE_DN" env-default:"dc=example,dc=com"`
	LDAPUserFilter string `env:"LDAP_USER_FILTER" env-default:"(objectClass=inetOrgPerson)"`
}

func NewLoadConfig() (Config, error) {
	var cfg Config
	cleanenv.ReadConfig(".env", &cfg)
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
