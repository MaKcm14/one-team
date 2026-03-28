package auth

import (
	"fmt"
	"math/rand"

	"github.com/MaKcm14/one-team/internal/services/usecase/user"
	"golang.org/x/crypto/bcrypt"
)

func (auth Interactor) checkPassword(origHashPwd string, password string, salt int) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(origHashPwd),
		[]byte(
			fmt.Sprintf("%d%s%d", salt, password, auth.cfg.GlobalPwdSalt),
		),
	)
}

func (auth Interactor) hashPassword(pwd string, userSalt int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(
		[]byte(
			fmt.Sprintf("%d%s%d", userSalt, pwd, auth.cfg.GlobalPwdSalt),
		),
		bcrypt.DefaultCost,
	)
}

func (auth Interactor) generateSalt() int {
	return rand.Intn(auth.cfg.GlobalPwdSalt)
}

func (auth Interactor) verifyPassword(pwd string) error {
	requiredSymbols := map[string]struct{}{
		"@": {},
		"!": {},
		"_": {},
		"?": {},
		"#": {},
		"$": {},
	}

	if len(pwd) < 9 || len(pwd) > 16 {
		return fmt.Errorf("%w: %w", user.ErrVerifyPassword, user.ErrPasswordLength)
	}

	count := 0
	for _, elem := range pwd {
		if _, ok := requiredSymbols[string(elem)]; ok {
			count++
		}
		if count == 2 {
			break
		}
	}

	if count != 2 {
		return fmt.Errorf("%w: %w", user.ErrVerifyPassword, user.ErrPasswordSymbols)
	}
	return nil
}
