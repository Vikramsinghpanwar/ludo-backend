package auth

import "time"

type User struct {
	ID           string
	Phone        string
	PasswordHash string
	Status       string
	CreatedAt    time.Time
}

type AuthOTP struct {
	ID                string
	Phone             string
	OTPHash           string
	PURPOSE           string
	OTP_Verify_Count  int //attempts to verify after otp sent
	OTP_Request_Count int //number of otp sent by this phone, for rate limiting
	Verified          bool
	LastSentAt        time.Time
	ExpiresAt         time.Time
}

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
}
