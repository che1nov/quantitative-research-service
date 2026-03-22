package usecases

import (
	"context"
	"errors"

	"github.com/che1nov/quantitative-research-service/internal/domain"
	"github.com/che1nov/quantitative-research-service/internal/dto"
)

// AuthenticateVKUseCase выполняет вход пользователя по VK ID.
type AuthenticateVKUseCase struct {
	users  UserStorage
	tokens TokenGenerator
	signer JWTSigner
	log    AppLogger
}

func NewAuthenticateVKUseCase(users UserStorage, tokens TokenGenerator, signer JWTSigner, log AppLogger) *AuthenticateVKUseCase {
	return &AuthenticateVKUseCase{users: users, tokens: tokens, signer: signer, log: log}
}

// Execute находит или создает пользователя и выдает JWT.
func (uc *AuthenticateVKUseCase) Execute(ctx context.Context, input dto.AuthenticateVKInput) (dto.AuthenticateVKOutput, error) {
	uc.log.InfoContext(ctx, "Аутентификация через VK", "vk_id", input.VKID)

	if input.VKID <= 0 {
		uc.log.WarnContext(ctx, "Получен некорректный vk id")
		return dto.AuthenticateVKOutput{}, domain.ErrInvalidVKID
	}

	user, err := uc.users.GetUserByVKID(ctx, input.VKID)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			uc.log.ErrorContext(ctx, "Ошибка чтения пользователя", "error", err)
			return dto.AuthenticateVKOutput{}, err
		}

		userID, err := uc.tokens.Generate(ctx, 12)
		if err != nil {
			uc.log.ErrorContext(ctx, "Ошибка генерации user id", "error", err)
			return dto.AuthenticateVKOutput{}, err
		}

		user, err = domain.NewVKUser(userID, input.VKID)
		if err != nil {
			return dto.AuthenticateVKOutput{}, err
		}
		if err := uc.users.CreateUser(ctx, user); err != nil {
			uc.log.ErrorContext(ctx, "Ошибка создания пользователя", "error", err)
			return dto.AuthenticateVKOutput{}, err
		}
		uc.log.InfoContext(ctx, "Создан новый пользователь", "user_id", user.ID)
	}

	token, err := uc.signer.SignUserToken(ctx, user)
	if err != nil {
		uc.log.ErrorContext(ctx, "Ошибка подписи JWT", "error", err)
		return dto.AuthenticateVKOutput{}, err
	}

	uc.log.InfoContext(ctx, "Аутентификация через VK выполнена", "user_id", user.ID)
	return dto.AuthenticateVKOutput{Token: token, UserID: user.ID, VKID: user.VKID}, nil
}
