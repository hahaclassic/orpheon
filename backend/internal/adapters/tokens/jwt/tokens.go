package jwttokens

import (
	"errors"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrExpired              = errors.New("access token expired error")
	ErrInvalidClaims        = errors.New("invalid claims error")
	ErrInvalidSigningMethod = errors.New("invalid signing method error")
	ErrTokenParsing         = errors.New("token parsing error")
	ErrInvalidAccessToken   = errors.New("invalid access token")
)

type jwtClaims struct {
	UserID    uuid.UUID          `json:"user_id"`
	AccessLvl entity.AccessLevel `json:"access_level"`
	jwt.RegisteredClaims
}

type AccessTokenConfig struct {
	TTL    time.Duration
	Jitter time.Duration
}

type JWTTokenService struct {
	secretKey []byte
	accessCnf AccessTokenConfig
}

func New(config AccessTokenConfig, secretKey []byte) *JWTTokenService {
	key := make([]byte, len(secretKey))
	copy(key, secretKey)

	return &JWTTokenService{
		accessCnf: config,
		secretKey: key,
	}
}

func (s *JWTTokenService) CreateAccessToken(claims *entity.Claims) (string, error) {
	jitter := time.Duration(rand.Int63n(int64(s.accessCnf.Jitter)))
	exp := time.Now().Add(s.accessCnf.TTL + jitter)

	jwtClaims := jwtClaims{
		UserID:    claims.UserID,
		AccessLvl: claims.AccessLvl,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	return token.SignedString(s.secretKey)
}

func (s *JWTTokenService) ParseAccessToken(tokenStr string) (*entity.Claims, error) {
	var jwtClaims jwtClaims

	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return s.secretKey, nil
	})

	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return nil, ErrExpired
	case errors.Is(err, jwt.ErrTokenInvalidClaims):
		return nil, ErrInvalidClaims
	case err != nil:
		return nil, ErrTokenParsing
	case !token.Valid:
		return nil, ErrInvalidAccessToken
	}

	return &entity.Claims{
		UserID:    jwtClaims.UserID,
		AccessLvl: jwtClaims.AccessLvl,
	}, nil
}
