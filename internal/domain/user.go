package domain

import "time"

// User описывает внутреннего пользователя приложения.
type User struct {
	ID        string
	VKID      int64
	CreatedAt time.Time
}

// NewVKUser создает пользователя по идентификатору VK.
func NewVKUser(id string, vkID int64) (User, error) {
	if vkID <= 0 {
		return User{}, ErrInvalidVKID
	}

	return User{
		ID:        id,
		VKID:      vkID,
		CreatedAt: time.Now().UTC(),
	}, nil
}
