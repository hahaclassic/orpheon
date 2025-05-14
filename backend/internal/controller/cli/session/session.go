package session

import (
	"sync"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	once sync.Once
	s    *Session
)

type Session struct {
	claims *entity.Claims
	tokens *entity.AuthTokens
	mu     *sync.RWMutex
}

func Instance() *Session {
	once.Do(func() {
		s = &Session{
			mu:     &sync.RWMutex{},
			claims: &entity.Claims{},
			tokens: &entity.AuthTokens{},
		}
	})
	return s
}

func StartSession(claims *entity.Claims, tokens *entity.AuthTokens) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.claims = claims
	s.tokens = tokens
}

func EndSession() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.claims.UserID = uuid.Nil
	s.claims.AccessLvl = entity.Unauthorized
	s.tokens.Access = ""
	s.tokens.Refresh = ""
}

func Claims() *entity.Claims {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.claims
}

func Tokens() *entity.AuthTokens {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.tokens
}

func UpdateTokens(tokens *entity.AuthTokens) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tokens = tokens
}

func IsAuthenticated() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.claims.UserID != uuid.Nil
}
