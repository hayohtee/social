package repository

import (
	"context"
	"crypto/sha256"
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
		case strings.Contains(err.Error(), `"users_email_key"`):
			return ErrDuplicateEmail
		case strings.Contains(err.Error(), `"users_username_key"`):
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

func (u *UsersRepository) CreateAndInvite(ctx context.Context, user *data.User, token []byte, invitationExp time.Duration) error {
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

func (u *UsersRepository) Activate(ctx context.Context, plainToken string) error {
	return withTx(u.db, ctx, func(tx *sql.Tx) error {
		user, err := getUserFromInvitation(ctx, tx, plainToken)
		if err != nil {
			return err
		}
		user.Activated = true
		if err := u.update(ctx, user, tx); err != nil {
			return err
		}

		if err := u.deleteUserInvitations(ctx, user.ID, tx); err != nil {
			return err
		}

		return nil
	})
}

func getUserFromInvitation(ctx context.Context, tx *sql.Tx, plainToken string) (data.User, error) {
	tokenHash := sha256.Sum256([]byte(plainToken))

	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.activated
		FROM users u
		INNER JOIN user_invitations ui ON u.id = ui.user_id	
		WHERE ui.token = $1 AND ui.expiry > $2`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	var user data.User
	err := tx.QueryRowContext(ctx, query, tokenHash[:], time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.Activated,
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

func (u *UsersRepository) update(ctx context.Context, user data.User, tx *sql.Tx) error {
	query := `
		UPDATE users SET username = $1, email = $2, activated = $3
		WHERE id = $4`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.Activated, user.ID)
	return err
}

func (u *UsersRepository) deleteUserInvitations(ctx context.Context, userID int64, tx *sql.Tx) error {
	query := `DELETE FROM user_invitations WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	return err
}
