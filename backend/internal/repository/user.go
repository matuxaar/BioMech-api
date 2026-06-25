package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/motvii/desertacia/internal/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, id, email string) (*model.User, error) {
	user := &model.User{
		ID:        id,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := r.db.Exec(ctx,
		`INSERT INTO users (id, email, created_at, updated_at)
		 VALUES ($1, $2, $3, $4)`,
		user.ID, user.Email, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateEmail(ctx context.Context, id, email string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET email = $1, updated_at = $2 WHERE id = $3`,
		email, time.Now(), id,
	)
	return err
}
