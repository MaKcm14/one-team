package repo

import "sync/atomic"

type repoCounters struct {
	userID    atomic.Int32
	accountID atomic.Int32
}
