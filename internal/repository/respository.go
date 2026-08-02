package repository

import (
	"context"
	"fmt"

	"automatomic/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	UpsertGitHubUser(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
}

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) UpsertGitHubUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (github_id, username, email, avatar_url, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (github_id) DO UPDATE SET
			username = EXCLUDED.username,
			email = EXCLUDED.email,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = NOW()
		RETURNING id, created_at, updated_at;
	`
	err := r.db.QueryRow(
		ctx, query,
		user.GitHubID, user.Username, user.Email, user.AvatarURL, user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("upsert github user: %w", err)
	}
	return nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, github_id, username, email, avatar_url, role, created_at, updated_at
		FROM users WHERE id = $1;
	`
	var user model.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.GitHubID, &user.Username, &user.Email,
		&user.AvatarURL, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}