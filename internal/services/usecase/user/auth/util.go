package auth

import (
	"fmt"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

func (auth Interactor) checkPassword(origHashPwd string, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(origHashPwd),
		[]byte(
			fmt.Sprintf("%s%d", password, auth.cfg.GlobalPwdSalt),
		),
	)
}

func (auth Interactor) hashPassword(pwd string, userSalt int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(
		[]byte(fmt.Sprintf("%s%d", pwd, auth.cfg.GlobalPwdSalt)),
		userSalt,
	)
}

func (auth Interactor) generateSalt() int {
	return rand.Intn(auth.cfg.GlobalPwdSalt)
}

func (auth Interactor) verifyPassword(pwd string) error {
	/*
		Requirments for the password:
		- length more than 9 symbols;
		- without any most use keys from the passwords dicts.
		- must be at the pwd: @, !, _, ?, #, $. // may be two of them.
	*/
	return nil
}
