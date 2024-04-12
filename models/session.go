package models

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/ThisJohan/go-htmx-chat/rand"
	"github.com/redis/go-redis/v9"
)

var (
	DefaultSessionBytes = 32
)

type SessionService struct {
	Redis *redis.Client
}

type UserCache struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (s *SessionService) Create(ctx context.Context, userData UserCache) (string, error) {
	sessionToken, err := rand.String(DefaultSessionBytes)
	if err != nil {
		return "", err
	}
	sessionTokenHash := s.hash(sessionToken)

	jsonData, _ := json.Marshal(userData)
	err = s.Redis.Set(ctx, sessionTokenHash, jsonData, time.Hour*1).Err()
	if err != nil {
		return "", err
	}
	return sessionToken, nil
}

func (s *SessionService) Get(ctx context.Context, sessionToken string) (*UserCache, error) {
	sessionTokenHash := s.hash(sessionToken)
	jsonData, err := s.Redis.Get(ctx, sessionTokenHash).Result()
	if err != nil {
		return nil, err
	}
	var userData UserCache
	err = json.Unmarshal([]byte(jsonData), &userData)
	if err != nil {
		return nil, err
	}
	return &userData, nil
}

func (*SessionService) hash(sessionToken string) string {
	hash := sha256.Sum256([]byte(sessionToken))
	return base64.URLEncoding.EncodeToString(hash[:])
}
