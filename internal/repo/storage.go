package repo

import (
	"auth-train/test/internal/entity"
	"log/slog"
)

type BankRepository struct {
	users    map[entity.UserID]entity.User
	counters repoCounters

	logger *slog.Logger
}

func NewBankRepository(logger *slog.Logger) BankRepository {
	return BankRepository{
		users:  make(map[entity.UserID]entity.User, 100),
		logger: logger,
	}
}

func (b *BankRepository) CreateUser(userCfg UserConfig) entity.User {
	user := UserConfigToUser(userCfg)

	b.counters.userID++
	user.ID = entity.UserID(b.counters.userID)

	b.counters.accountID++
	user.Account.ID = entity.AccountID(b.counters.accountID)

	b.users[user.ID] = user
	return user
}

func (b *BankRepository) DeleteUser(id entity.UserID) {
	delete(b.users, id)
}

func (b *BankRepository) GetUsers() []entity.User {
	users := make([]entity.User, 0, len(b.users))
	for _, user := range b.users {
		users = append(users, user)
	}
	return users
}

func (b *BankRepository) GetUser(id entity.UserID) (entity.User, error) {
	user, ok := b.users[id]
	if !ok {
		return entity.User{}, ErrUserNotExist
	}
	return user, nil
}

func (b *BankRepository) SetMoney(id entity.UserID, money float64) (entity.User, error) {
	user, ok := b.users[id]
	if !ok {
		b.logger.Warn(ErrUserNotExist.Error())
		return entity.User{}, ErrUserNotExist
	}
	user.Account.Money = money
	b.users[id] = user

	return user, nil
}
