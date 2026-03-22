package token

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/MaKcm14/one-team/internal/config"
	"golang.org/x/crypto/bcrypt"
)

const RefreshTokenTTL = 24 * time.Hour

type RefreshToken struct {
	cfg config.AuthConfig
}

func NewRefreshToken(cfg config.AuthConfig) RefreshToken {
	return RefreshToken{
		cfg: cfg,
	}
}

func genOpaqueToken(size int) string {
	blockSize := size / 12

	buff := make([]byte, size)
	for i := 0; i < size; {
		if i+blockSize >= size {
			genRandSlice(buff[i:])
		} else {
			genRandSlice(buff[i : i+blockSize])
		}
		i += blockSize
	}
	return string(buff)
}

func genRandSlice(slice []byte) {
	for i := 0; i != len(slice); i++ {
		num, caseForm := rand.Intn(15), rand.Intn(2)

		if caseForm == 0 {
			slice[i] = []byte(fmt.Sprintf("%x", num))[0]
		} else {
			slice[i] = []byte(fmt.Sprintf("%X", num))[0]
		}
	}
}

func (r RefreshToken) IssueRefreshToken(size int) string {
	return genOpaqueToken(size)
}

func (r RefreshToken) CheckRefreshToken(origHashedToken string, token string) error {
	return bcrypt.CompareHashAndPassword([]byte(origHashedToken), []byte(token))
}

func (r RefreshToken) HashRefreshToken(token string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
}
