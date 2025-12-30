package auth

import (
	"context"
	"database/sql"
	"errors"
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
