package repo

import "auth-train/test/internal/entity"

type UserConfig struct {
	Name        string `json:"name,omitempty"`
	Surname     string `json:"surname,omitempty"`
	Passport    string `json:"passport,omitempty"`
	AdminStatus bool   `json:"admin,omitempty"`
	PwdHash     []byte `json:"-"`
}

func UserConfigToUser(userCfg UserConfig) entity.User {
	return entity.User{
		Name:     userCfg.Name,
		Surname:  userCfg.Surname,
		Passport: userCfg.Passport,
		Profile: entity.UserProfile{
			PwdHash:     userCfg.PwdHash,
			AdminStatus: userCfg.AdminStatus,
		},
	}
}
