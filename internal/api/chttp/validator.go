package chttp

import (
	"fmt"
	"regexp"
	"strconv"

	"auth-train/test/internal/entity"
	"auth-train/test/internal/repo"
)

const (
	userIDParamName = "userID"
)

func validateUserID(userIDParam string) (entity.UserID, error) {
	id, err := strconv.Atoi(userIDParam)
	if err != nil {
		return 0, ErrRequestQueryParam
	}
	return entity.UserID(id), nil
}

type userOpt func(*repo.UserConfig) error

func withNameValidator() userOpt {
	return func(userCfg *repo.UserConfig) error {
		if len(userCfg.Name) == 0 {
			return fmt.Errorf("%w: error of the user's Name", ErrInvalidValue)
		}
		return nil
	}
}

func withSurnameValidator() userOpt {
	return func(userCfg *repo.UserConfig) error {
		if len(userCfg.Surname) == 0 {
			return fmt.Errorf("%w: error of the user's Surname", ErrInvalidValue)
		}
		return nil
	}
}

func withPassportValidator() userOpt {
	return func(userCfg *repo.UserConfig) error {
		passReg := regexp.MustCompile(`^\d{4} \d{6}$`)
		if !passReg.MatchString(userCfg.Passport) {
			return fmt.Errorf("%w: error of the user's Passport Format",
				ErrInvalidValue)
		}
		return nil
	}
}

func validateUser(userCfg repo.UserConfig, validators ...userOpt) error {
	for _, opt := range validators {
		if err := opt(&userCfg); err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidData, err)
		}
	}
	return nil
}
