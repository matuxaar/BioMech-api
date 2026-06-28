package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matuxaar/BioMech-api/internal/model"
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
		`SELECT id, email, COALESCE(nickname, ''), COALESCE(display_name, ''), COALESCE(photo_url, ''), created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&user.ID, &user.Email, &user.Nickname, &user.DisplayName, &user.PhotoURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByNickname(ctx context.Context, nickname string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, COALESCE(nickname, ''), COALESCE(display_name, ''), COALESCE(photo_url, ''), created_at, updated_at
		 FROM users WHERE nickname = $1`, nickname,
	).Scan(&user.ID, &user.Email, &user.Nickname, &user.DisplayName, &user.PhotoURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) CountDevices(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM devices WHERE user_id = $1`, userID,
	).Scan(&count)
	return count, err
}

func (r *UserRepository) UpdateEmail(ctx context.Context, id, email string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET email = $1, updated_at = $2 WHERE id = $3`,
		email, time.Now(), id,
	)
	return err
}

func (r *UserRepository) UpdateProfile(ctx context.Context, id string, req *model.UpdateUserRequest) (*model.User, error) {
	if req.Nickname != nil {
		_, err := r.db.Exec(ctx, `UPDATE users SET nickname = $1, updated_at = $2 WHERE id = $3`, *req.Nickname, time.Now(), id)
		if err != nil {
			return nil, err
		}
	}
	if req.DisplayName != nil {
		_, err := r.db.Exec(ctx, `UPDATE users SET display_name = $1, updated_at = $2 WHERE id = $3`, *req.DisplayName, time.Now(), id)
		if err != nil {
			return nil, err
		}
	}
	if req.PhotoURL != nil {
		_, err := r.db.Exec(ctx, `UPDATE users SET photo_url = $1, updated_at = $2 WHERE id = $3`, *req.PhotoURL, time.Now(), id)
		if err != nil {
			return nil, err
		}
	}
	return r.FindByID(ctx, id)
}
