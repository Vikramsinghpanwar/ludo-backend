package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type PostgresAuthRepo struct {
	db *sql.DB
}

func NewPostgresAuthRepo(db *sql.DB) *PostgresAuthRepo {
	return &PostgresAuthRepo{db: db}
}

func (r *PostgresAuthRepo) CreateUser(ctx context.Context, user *User) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (phone, password_hash, status)
		 VALUES ($1, $2, $3)`,
		user.Phone,
		user.PasswordHash,
		user.Status,
	)
	return err
}

func (r *PostgresAuthRepo) GetUserByPhone(ctx context.Context, phone string) (*User, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, phone, password_hash, status
		 FROM users WHERE phone = $1`,
		phone,
	)

	u := &User{}
	err := row.Scan(&u.ID, &u.Phone, &u.PasswordHash, &u.Status)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *PostgresAuthRepo) SaveRefreshToken(ctx context.Context, token *RefreshToken) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		 VALUES ($1, $2, $3)`,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
	)
	return err
}

func (r *PostgresAuthRepo) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM refresh_tokens WHERE token_hash = $1`,
		tokenHash,
	)
	return err
}

func (r *PostgresAuthRepo) GetRefreshToken(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT user_id, token_hash, expires_at
		 FROM refresh_tokens WHERE token_hash = $1`,
		tokenHash,
	)

	rt := &RefreshToken{}
	err := row.Scan(&rt.UserID, &rt.TokenHash, &rt.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("refresh token not found")
	}
	if err != nil {
		return nil, err
	}

	return rt, nil
}

func (r *PostgresAuthRepo) SaveOTPRequest(ctx context.Context, AuthOTP *AuthOTP) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO auth_otp_requests (phone, otp_hash, purpose, attempts, verify_attempts, otp_verified, last_sent_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (phone)
		DO UPDATE SET
			otp_hash = EXCLUDED.otp_hash,
			purpose = EXCLUDED.purpose,
			expires_at = EXCLUDED.expires_at,
			attempts = EXCLUDED.attempts,
			verify_attempts = EXCLUDED.verify_attempts,
			otp_verified = EXCLUDED.otp_verified,
			last_sent_at =EXCLUDED.last_sent_at`,
		AuthOTP.Phone,
		AuthOTP.OTPHash,
		AuthOTP.PURPOSE,
		AuthOTP.OTP_Request_Count,
		AuthOTP.OTP_Verify_Count,
		AuthOTP.Verified,
		AuthOTP.LastSentAt,
		AuthOTP.ExpiresAt,
	)
	return err
}

func (r *PostgresAuthRepo) GetOTPRequestByPhone(ctx context.Context, phone string) (*AuthOTP, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, phone, otp_hash, attempts, verify_attempts, otp_verified, last_sent_at, expires_at
		 FROM auth_otp_requests WHERE phone = $1`,
		phone,
	)

	tu := &AuthOTP{}
	err := row.Scan(&tu.ID, &tu.Phone, &tu.OTPHash, &tu.OTP_Request_Count, &tu.OTP_Verify_Count, &tu.Verified, &tu.LastSentAt, &tu.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("otp Req not found")
	}
	if err != nil {
		return nil, err
	}

	return tu, nil
}

func (r *PostgresAuthRepo) DeleteOTPRequestByPhone(ctx context.Context, phone string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM auth_otp_requests WHERE phone = $1`,
		phone,
	)
	return err
}

func (r *PostgresAuthRepo) VerifyOTPTransaction(
	ctx context.Context,
	phone string,
	otpHash string,
	maxAttempts int,
) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // safe no-op if committed

	var (
		dbOTPHash      string
		verifyAttempts int
		expiresAt      time.Time
		otpVerified    bool
	)

	// 🔒 Lock row so no parallel verification can happen
	err = tx.QueryRowContext(
		ctx,
		`
		SELECT otp_hash, verify_attempts, expires_at, otp_verified
		FROM auth_otp_requests
		WHERE phone = $1
		FOR UPDATE
		`,
		phone,
	).Scan(&dbOTPHash, &verifyAttempts, &expiresAt, &otpVerified)

	if err == sql.ErrNoRows {
		return errors.New("otp not found")
	}
	if err != nil {
		return err
	}

	// Already verified
	if otpVerified {
		return errors.New("otp already verified")
	}

	// Expired
	if time.Now().After(expiresAt) {
		return errors.New("expired otp")
	}

	// Too many attempts
	if verifyAttempts >= maxAttempts {
		return errors.New("too many otp attempts")
	}

	// Wrong OTP → increment attempts
	if dbOTPHash != otpHash {
		_, err = tx.ExecContext(
			ctx,
			`
			UPDATE auth_otp_requests
			SET verify_attempts = verify_attempts + 1
			WHERE phone = $1
			`,
			phone,
		)
		if err != nil {
			return err
		}

		return errors.New("invalid otp")
	}

	// ✅ Correct OTP → mark verified
	_, err = tx.ExecContext(
		ctx,
		`
		UPDATE auth_otp_requests
		SET otp_verified = true
		WHERE phone = $1
		`,
		phone,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}
