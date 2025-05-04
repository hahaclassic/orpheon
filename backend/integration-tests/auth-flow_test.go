package integration_test

// Я хочу протестировать сценарий регистрации-входа-выхода пользователя. В этом у меня участвует несколько репозиториев:
// type RefreshTokenRepository interface {
// 	Set(ctx context.Context, token string, claims *entity.Claims) error
// 	Get(ctx context.Context, token string) (*entity.Claims, error)
// 	Delete(ctx context.Context, token string) error
// }

// type AuthRepository interface {
// 	SaveCredentials(ctx context.Context, userID uuid.UUID, credentials *entity.UserCredentials) error
// 	GetPasswordByLogin(ctx context.Context, login string) (string, error)
// 	GetPasswordByID(ctx context.Context, userID uuid.UUID) (string, error)
// 	GetClaimsByLogin(ctx context.Context, login string) (*entity.Claims, error)
// 	UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error
// }

// type UserCreatorService interface {
// 	CreateUser(ctx context.Context, info *entity.UserInfo) (uuid.UUID, error)
// }

// По сути, при регистрации мне надо:
// 1. создать юзера (в постгрес)
// 2. сохранить его логин/пароль  (в постгрес)
// 3. автоматический вход, поэтому сохраняется рефреш токен (в редисе)
