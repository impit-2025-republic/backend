package service

import (
	"b8boost/backend/internal/entities"
	"b8boost/backend/internal/infra/ldap"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type LDAPService struct {
	ldap       ldap.LDAP
	userRepo   entities.UserRepo
	userWallet entities.UserWalletRepo
}

func NewLDAPService(
	ldap ldap.LDAP,
	userRepo entities.UserRepo,
	userWallet entities.UserWalletRepo,
) LDAPService {
	return LDAPService{
		ldap:       ldap,
		userRepo:   userRepo,
		userWallet: userWallet,
	}
}

func (s LDAPService) Sync() {
	ldapUsers, err := s.ldap.FetchAllUsers()

	if err != nil {
		fmt.Println(err)
		return
	}

	dbUsers, err := s.userRepo.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	ldapUserMap := make(map[string]ldap.LDAPUserData)
	dbUserMap := make(map[string]entities.User)

	for _, ldapUser := range ldapUsers {
		uid := ldap.GetFirstValueOrDefault(ldapUser, "entryUUID", "")
		if uid != "" {
			ldapUserMap[uid] = ldapUser
		}
	}

	for i := range dbUsers {
		dbUserMap[dbUsers[i].LdapID] = dbUsers[i]
	}

	for uid, ldapUser := range ldapUserMap {
		if dbUser, exists := dbUserMap[uid]; exists {
			if err := s.updateUser(dbUser, ldapUser); err != nil {
				log.Printf("Error updating user %s: %v", uid, err)
				continue
			}
			log.Printf("Updated user: %s", uid)
		} else {
			if err := s.createUser(ldapUser); err != nil {
				log.Printf("Error creating user %s: %v", uid, err)
				continue
			}
			log.Printf("Created new user: %s", uid)
		}
	}

	// for uid, dbUser := range dbUserMap {
	// 	if _, exists := ldapUserMap[uid]; !exists && dbUser.Active {
	// 		if err := sm.deactivateUser(ctx, dbUser); err != nil {
	// 			log.Printf("Error deactivating user %s: %v", uid, err)
	// 			continue
	// 		}
	// 		log.Printf("Deactivated user: %s", uid)
	// 	}
	// }
}

func (s LDAPService) createUser(ldapUser ldap.LDAPUserData) error {
	var dbUser entities.User
	tgDef := 0
	if dbUser.TelegramID != nil {
		tgDef = *dbUser.TelegramID
	}

	layout := "20060102150405Z"
	createdAt, err := time.Parse(layout, ldap.GetFirstValueOrDefault(ldapUser, "createTimestamp", ""))
	if err != nil {
		createdAt = time.Now()
	}

	updatedAt, err := time.Parse(layout, ldap.GetFirstValueOrDefault(ldapUser, "modifyTimestamp", ""))
	if err != nil {
		updatedAt = time.Now()
	}

	tgId := ldap.GetFirstValueOrDefaultInt(ldapUser, "description", tgDef)

	dbUser.Surname = ldap.GetFirstValueOrDefaultPtr(ldapUser, "sn", dbUser.Surname)
	dbUser.Email = ldap.GetFirstValueOrDefaultPtr(ldapUser, "mail", dbUser.Email)
	dbUser.TelegramID = &tgId
	dbUser.Name = ldap.GetFirstValueOrDefaultPtr(ldapUser, "cn", dbUser.Name)
	dbUser.LastSurname = ldap.GetFirstValueOrDefaultPtr(ldapUser, "givenName", dbUser.LastSurname)
	dbUser.CreatedAt = createdAt
	dbUser.UpdatedAt = &updatedAt
	dbUser.Phone = ldap.GetFirstValueOrDefaultPtr(ldapUser, "mobile", dbUser.Phone)
	dbUser.LdapID = ldap.GetFirstValueOrDefault(ldapUser, "entryUUID", "")

	user, err := s.userRepo.Create(dbUser)
	if err != nil {
		return err
	}
	if tgId != 0 {
		s.userWallet.Create(entities.UserWallet{
			UserID: int(user.UserID),
			Price:  0.0,
		})
	}
	return nil
}

func (s LDAPService) updateUser(dbUser entities.User, ldapUser ldap.LDAPUserData) error {
	tgDef := 0
	if dbUser.TelegramID != nil {
		tgDef = *dbUser.TelegramID
	}

	layout := "20060102150405Z"
	createdAt, err := time.Parse(layout, ldap.GetFirstValueOrDefault(ldapUser, "createTimestamp", ""))
	if err != nil {
		createdAt = time.Now()
	}

	updatedAt, err := time.Parse(layout, ldap.GetFirstValueOrDefault(ldapUser, "modifyTimestamp", ""))
	if err != nil {
		updatedAt = time.Now()
	}

	tgId := ldap.GetFirstValueOrDefaultInt(ldapUser, "description", tgDef)
	dbUser.Surname = ldap.GetFirstValueOrDefaultPtr(ldapUser, "sn", dbUser.Surname)
	dbUser.Email = ldap.GetFirstValueOrDefaultPtr(ldapUser, "mail", dbUser.Email)
	dbUser.TelegramID = &tgId
	dbUser.Name = ldap.GetFirstValueOrDefaultPtr(ldapUser, "cn", dbUser.Name)
	dbUser.LastSurname = ldap.GetFirstValueOrDefaultPtr(ldapUser, "givenName", dbUser.LastSurname)
	dbUser.CreatedAt = createdAt
	dbUser.UpdatedAt = &updatedAt
	dbUser.Phone = ldap.GetFirstValueOrDefaultPtr(ldapUser, "mobile", dbUser.Phone)
	dbUser.LdapID = ldap.GetFirstValueOrDefault(ldapUser, "entryUUID", "")

	if tgId != 0 {
		_, err := s.userWallet.GetWallet(dbUser.UserID)
		if err != nil && err == gorm.ErrRecordNotFound {
			s.userWallet.Create(entities.UserWallet{
				UserID: int(dbUser.UserID),
				Price:  0.0,
			})
		}
	}
	return s.userRepo.Update(dbUser)
}

// func (s LDAPService) deactivateUser(ctx context.Context, dbUser *entities.User) error {
// 	dbUser. = time.Now()

// 	return s.userRepo.Update(dbUser)
// }
