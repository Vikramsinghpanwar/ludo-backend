package auth

import (
	"context"
	"errors"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Signup(ctx context.Context, req SignupRequest) error {
	if len(req.Password) < 6 {
		return errors.New("weak password")
	}

	hash := hashString(req.Password)

	user := &User{
		Phone:        req.Phone,
		PasswordHash: hash,
		Status:       "active",
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.GetUserByPhone(ctx, req.Phone)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if hashString(req.Password) != user.PasswordHash {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := generateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshPlain, refreshHash, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	rt := &RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	err = s.repo.SaveRefreshToken(ctx, rt)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshPlain,
	}, nil
}

func (s *Service) Logout(ctx context.Context, tokenHash string) error {
	return s.repo.DeleteRefreshToken(ctx, tokenHash)
}

func (s *Service) DeleteRefreshToken(ctx context.Context, tonekHash string) error {
	return s.repo.DeleteRefreshToken(ctx, tonekHash)
}

func (s *Service) Refresh(ctx context.Context, refreshPlain string) (*AuthResponse, error) {
	hash := hashString(refreshPlain)

	rt, err := s.repo.GetRefreshToken(ctx, hash)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// rotate
	s.repo.DeleteRefreshToken(ctx, hash)

	access, _ := generateAccessToken(rt.UserID)
	newPlain, newHash, _ := generateRefreshToken()

	s.repo.SaveRefreshToken(ctx, &RefreshToken{
		UserID:    rt.UserID,
		TokenHash: newHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})

	return &AuthResponse{
		AccessToken:  access,
		RefreshToken: newPlain,
	}, nil
}
