package models

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ThisJohan/go-htmx-chat/rand"
	"github.com/redis/go-redis/v9"
)

var (
	DefaultSessionBytes = 32
)

type SessionService struct {
	Redis *redis.Client
}

func (s *SessionService) Create(ctx context.Context, userData interface{}) (string, error) {
	sessionToken, err := rand.String(DefaultSessionBytes)
	if err != nil {
		return "", err
	}
	sessionTokenHash := s.hash(sessionToken)

	userData, _ = json.Marshal(userData)
	err = s.Redis.Set(ctx, sessionTokenHash, userData, 0).Err()
	if err != nil {
		return "", err
	}
	return sessionToken, nil
}

func (s *SessionService) Test(ctx context.Context) {
	fmt.Println("Test")
}

func (*SessionService) hash(sessionToken string) string {
	hash := sha256.Sum256([]byte(sessionToken))
	return base64.URLEncoding.EncodeToString(hash[:])
}
