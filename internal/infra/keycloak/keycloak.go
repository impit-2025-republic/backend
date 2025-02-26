package keycloak

import (
	"b8boost/backend/config"
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

type Keycloak struct {
	server        string
	realm         string
	adminUsername string
	adminPassword string
}

func NewKeycloak(conf config.Config) Keycloak {
	return Keycloak{
		server:        conf.KeycloakServer,
		realm:         "",
		adminUsername: "",
		adminPassword: "",
	}
}

func (k Keycloak) connect(client *gocloak.GoCloak) (*gocloak.JWT, error) {
	ctx := context.Background()
	token, err := client.LoginAdmin(ctx, k.adminUsername, k.adminPassword, k.realm)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (k Keycloak) CreateUser(username string, email string, firstName string, lastName string, tgID *int64, phone string) error {
	client := gocloak.NewClient(k.server)
	token, err := k.connect(client)
	if err != nil {
		fmt.Println(err)
		return err
	}
	enabled := true
	emailVerify := true
	user := gocloak.User{
		Username:      &username,
		Email:         &email,
		FirstName:     &firstName,
		LastName:      &lastName,
		Enabled:       &enabled,
		EmailVerified: &emailVerify,
		Attributes: &map[string][]string{
			"phone": {phone},
		},
	}

	_, err = client.CreateUser(context.TODO(), token.AccessToken, k.realm, user)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
