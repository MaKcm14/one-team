package auth

import "golang.org/x/crypto/bcrypt"

func (auth Interactor) checkPassword(origHashPwd string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(origHashPwd), []byte(password))
}
