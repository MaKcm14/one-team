package repo

import "auth-train/test/internal/entity"

type UserConfig struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Passport string `json:"passport"`
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
