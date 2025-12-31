package auth

import "context"

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByPhone(ctx context.Context, phone string) (*User, error)
	SaveRefreshToken(ctx context.Context, token *RefreshToken) error
	GetRefreshToken(ctx context.Context, tokenHash string) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, tokenHash string) error

	SaveOTPRequest(ctx context.Context, AuthOTP *AuthOTP) error
	GetOTPRequestByPhone(ctx context.Context, phone string) (*AuthOTP, error)
	DeleteOTPRequestByPhone(ctx context.Context, phone string) error
	VerifyOTPTransaction(ctx context.Context, phone string, otpHash string, maxAttempts int) error
}
