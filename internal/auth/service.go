package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/vikramsinghpanwar/ludo-backend/internal/infra/sms"
)

type Service struct {
	repo      Repository
	smsSender sms.Sender
}

func NewService(r Repository, sender sms.Sender) *Service {
	return &Service{
		repo:      r,
		smsSender: sender,
	}
}

const maxOTPVerifications int = 3
const maxOTPRequest int = 5
const OTPResendColldown = 30 * time.Second

func (s *Service) Signup(ctx context.Context, req SignupRequest) error {
	if len(req.Password) < 6 {
		return errors.New("password must be 6 digit long")
	}

	hash := hashString(req.Password)

	user := &User{
		Phone:        req.Phone,
		PasswordHash: hash,
		Status:       "active",
	}

	tmpUser, err := s.repo.GetOTPRequestByPhone(ctx, req.Phone)
	if err != nil {
		return err
	}
	fmt.Println(tmpUser.Verified)

	if tmpUser.Verified == false {

		return errors.New("Not verified")
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *Service) OTP(ctx context.Context, req OTPRequest) error {

	//otp limiting logic

	otpReq, _ := s.repo.GetOTPRequestByPhone(ctx, req.Phone)

	if otpReq != nil {
		if otpReq.OTP_Request_Count >= maxOTPRequest {
			return errors.New("OTP Limit Reached")
		}

		if time.Since(otpReq.LastSentAt) <= OTPResendColldown {
			return errors.New("Try after few seconds")
		}
	}

	otp, error := generateOTP()

	if error != nil {
		return fmt.Errorf("failed to generate otp: %w", error)
	}
	err := s.smsSender.SendOTP(req.Phone, otp)
	if err != nil {
		return err
	}

	tmpUser := &AuthOTP{
		Phone:             req.Phone,
		OTPHash:           hashString(otp),
		PURPOSE:           req.Purpose,
		OTP_Request_Count: 1,
		OTP_Verify_Count:  0,
		Verified:          false,
		LastSentAt:        time.Now(),
		ExpiresAt:         time.Now().Add(15 * time.Minute),
	}

	e := s.repo.SaveOTPRequest(ctx, tmpUser)
	if e != nil {
		return e
	}

	return nil
}
func (s *Service) VerifyOTP(ctx context.Context, req VerifyOTPRequest) (string, error) {
	err := s.repo.VerifyOTPTransaction(
		ctx,
		req.Phone,
		hashString(req.OTP),
		maxOTPVerifications,
	)
	if err != nil {
		return "", err
	}

	return "verification successful", nil
}

func generateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
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
