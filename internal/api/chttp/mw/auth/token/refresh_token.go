package token

import (
	"fmt"
	"math/rand"
	"time"
)

const RefreshTokenTTL = 24 * time.Hour

func genOpaqueToken(size int) string {
	blockSize := size / 12

	buff := make([]byte, size)
	for i := 0; i < size; {
		if i+blockSize >= size {
			genRandSlice(buff[i:])
		} else {
			genRandSlice(buff[i : i+blockSize])
		}
	}
	return string(buff)
}

func genRandSlice(slice []byte) {
	for i := 0; i != len(slice); i++ {
		num, caseForm := rand.Intn(15), rand.Intn(1)

		if caseForm == 0 {
			slice[i] = []byte(fmt.Sprintf("%x", num))[0]
		} else {
			slice[i] = []byte(fmt.Sprintf("%X", num))[0]
		}
	}
}

func IssueRefreshToken(size int) string {
	return genOpaqueToken(size)
}
