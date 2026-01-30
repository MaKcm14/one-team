package repo

import (
	"log/slog"

	"auth-train/test/internal/entity"
)

type Repository struct {
	bankUsers map[entity.UserID]entity.User
	counters  repoCounters

	logger *slog.Logger
}

func NewRepository(logger *slog.Logger) Repository {
	return Repository{
		bankUsers: make(map[entity.UserID]entity.User, 100),
		logger:    logger,
	}
}

func (b *Repository) CreateUser(userCfg UserConfig) entity.User {
	user := UserConfigToUser(userCfg)

	b.counters.userID++
	user.ID = entity.UserID(b.counters.userID)

	b.counters.accountID++
	user.Account.ID = entity.BankAccountID(b.counters.accountID)

	b.bankUsers[user.ID] = user
	return user
}

func (b *Repository) DeleteUser(id entity.UserID) {
	delete(b.bankUsers, id)
}

func (b *Repository) GetUsers() []entity.User {
	users := make([]entity.User, 0, len(b.bankUsers))
	for _, user := range b.bankUsers {
		users = append(users, user)
	}
	return users
}

func (b *Repository) GetUser(id entity.UserID) (entity.User, error) {
	user, ok := b.bankUsers[id]
	if !ok {
		return entity.User{}, ErrUserNotExist
	}
	return user, nil
}

func (b *Repository) GetUserByPassport(passport string) (entity.User, error) {
	for _, user := range b.bankUsers {
		if user.Passport == passport {
			return user, nil
		}
	}
	return entity.User{}, ErrUserNotExist
}

func (b *Repository) SetMoney(id entity.UserID, money float64) (entity.User, error) {
	user, ok := b.bankUsers[id]
	if !ok {
		return entity.User{}, ErrUserNotExist
	}
	user.Account.Money = money
	b.bankUsers[id] = user

	return user, nil
}

func (b *Repository) SetAdminStatus(id entity.UserID, status bool) error {
	user, ok := b.bankUsers[id]
	if !ok {
		return ErrUserNotExist
	}
	user.Profile.AdminStatus = status
	b.bankUsers[id] = user

	return nil
}
