package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ticketapp/internal/models"
)

/* =========================
   REPOSITORY
========================= */

type PostgresRefreshTokenRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRefreshTokenRepo(db *pgxpool.Pool) *PostgresRefreshTokenRepo {
	return &PostgresRefreshTokenRepo{db: db}
}

/* =========================
   STORE TOKEN
========================= */

func (r *PostgresRefreshTokenRepo) Store(
	userID uuid.UUID,
	hash string,
	exp time.Time,
) (uuid.UUID, error) {

	id := uuid.New()

	_, err := r.db.Exec(
		context.Background(),
		`INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
		 VALUES ($1,$2,$3,$4)`,
		id, userID, hash, exp,
	)

	return id, err
}

/* =========================
   GET VALID TOKEN
========================= */

func (r *PostgresRefreshTokenRepo) GetValid(hash string) (*models.RefreshToken, error) {
	t := &models.RefreshToken{}

	err := r.db.QueryRow(
		context.Background(),
		`SELECT 
			id,
			user_id,
			revoked,
			replaced_by,
			expires_at
		FROM refresh_tokens
		WHERE token_hash = $1`,
		hash,
	).Scan(
		&t.ID,
		&t.UserID,
		&t.Revoked,
		&t.ReplacedBy,
		&t.ExpiresAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	// ⏰ expired
	if t.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	// 🚨 REAL reuse detection
	if t.Revoked && t.ReplacedBy != nil {
		return nil, ErrTokenReuse
	}

	// ❌ normal revoked (logout / cleanup)
	if t.Revoked {
		return nil, ErrTokenRevoked
	}

	return t, nil
}

/* =========================
   MARK TOKEN REPLACED
========================= */

func (r *PostgresRefreshTokenRepo) MarkReplaced(
	oldID uuid.UUID,
	newID uuid.UUID,
) error {

	cmd, err := r.db.Exec(
		context.Background(),
		`UPDATE refresh_tokens
		 SET revoked = true, replaced_by = $2
		 WHERE id = $1`,
		oldID, newID,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("token not found")
	}

	return nil
}

/* =========================
   ROTATE TOKEN (CRITICAL)
========================= */

func (r *PostgresRefreshTokenRepo) Rotate(
	old *models.RefreshToken,
	newHash string,
	expiry time.Time,
) error {

	newID, err := r.Store(old.UserID, newHash, expiry)
	if err != nil {
		return err
	}

	return r.MarkReplaced(old.ID, newID)
}

/* =========================
   REVOKE
========================= */

func (r *PostgresRefreshTokenRepo) Revoke(tokenID uuid.UUID) error {
	cmd, err := r.db.Exec(
		context.Background(),
		`UPDATE refresh_tokens SET revoked=true WHERE id=$1`,
		tokenID,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("token not found")
	}

	return nil
}

func (r *PostgresRefreshTokenRepo) RevokeAll(userID uuid.UUID) error {
	_, err := r.db.Exec(
		context.Background(),
		`UPDATE refresh_tokens SET revoked=true WHERE user_id=$1`,
		userID,
	)
	return err
}
