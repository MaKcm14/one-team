package repo

import "auth-train/test/internal/entity"

type UserConfig struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Passport string `json:"passport"`
}

func UserConfigToUser(userCfg UserConfig) entity.User {
	user := entity.User{}
	user.Name = userCfg.Name
	user.Surname = userCfg.Surname
	user.Passport = userCfg.Passport
	return user
}
