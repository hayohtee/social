package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/hayohtee/social/internal/data"
)

type UsersRepository struct {
	db *sql.DB
}

func (u *UsersRepository) Create(ctx context.Context, user *data.User, tx *sql.Tx) error {
	query := `
		INSERT INTO users(username, email, password)
		VALUES($1, $2, $3)
		RETURNING id, created_at`

	args := []any{user.Username, user.Email, user.Password.Hash}

	if tx == nil {
		return u.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password.Hash).Scan(
			&user.ID,
			&user.CreatedAt,
		)
	}

	err := tx.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), `pq: duplicate key violates unique constraint "users_email_key"`):
			return ErrDuplicateEmail
		case strings.Contains(err.Error(), `pq: duplicate key violates unique constraint "users_username_key"`):
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (u *UsersRepository) GetByID(ctx context.Context, id int64) (data.User, error) {
	query := `
		SELECT id, username, email, created_at
		FROM users
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	var user data.User
	err := u.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return data.User{}, ErrNotFound
		default:
			return data.User{}, err
		}
	}
	return user, nil
}

func (u *UsersRepository) CreateAndInvite(ctx context.Context, user *data.User, token string, invitationExp time.Duration) error {
	return withTx(u.db, ctx, func(tx *sql.Tx) error {
		if err := u.Create(ctx, user, tx); err != nil {
			return err
		}
		return u.createInvitation(ctx, tx, token, invitationExp, user.ID)
	})
}

func (u *UsersRepository) createInvitation(ctx context.Context, tx *sql.Tx, token []byte, invitationExp time.Duration, userID int64) error {
	query := `
		INSERT INTO user_invitations(token, user_id, expiry)
		VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationExp))
	if err != nil {
		return err
	}
	return nil
}
