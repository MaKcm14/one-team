package repo

import "auth-train/test/internal/entity"

type UserConfig struct {
	Name     string `json:"name,omitempty"`
	Surname  string `json:"surname,omitempty"`
	Passport string `json:"passport,omitempty"`
	PwdHash  []byte `json:"-"`
}

func UserConfigToUser(userCfg UserConfig) entity.User {
	user := entity.User{}
	user.Name = userCfg.Name
	user.Surname = userCfg.Surname
	user.Passport = userCfg.Passport
	user.PwdHash = userCfg.PwdHash
	return user
}
