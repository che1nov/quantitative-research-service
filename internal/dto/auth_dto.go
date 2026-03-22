package dto

// AuthenticateVKInput описывает входные данные для входа через VK.
type AuthenticateVKInput struct {
	VKID int64 `json:"vk_id"`
}

// AuthenticateVKOutput описывает ответ после успешной аутентификации.
type AuthenticateVKOutput struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
	VKID   int64  `json:"vk_id"`
}
