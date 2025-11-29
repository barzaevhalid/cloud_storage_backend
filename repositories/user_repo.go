package repositories

import (
	"context"

	"github.com/barzaevhalid/cloud_storage_backend/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {

	return &UserRepository{
		db: db,
	}

}

func (r *UserRepository) Create(u *models.User) error {
	query := `INSERT INTO users (email, passwordhash, fullname) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRow(context.Background(), query, u.Email, u.PasswordHash, u.FullName).Scan(&u.ID)
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, email FROM users WHERE email=$1`

	err := r.db.QueryRow(context.Background(), query, email).Scan(&u.ID, &u.Email)

	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) Login(email string) (*models.User, error) {
	u := &models.User{}

	query := `SELECT id, email, fullname, passwordhash FROM users WHERE email=$1`

	err := r.db.QueryRow(context.Background(), query, email).Scan(&u.ID, &u.Email, &u.FullName, &u.PasswordHash)

	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) GetMe(id int64) (*models.User, error) {
	u := &models.User{}

	query := `SELECT id, email, fullname  FROM users WHERE id=$1`

	err := r.db.QueryRow(context.Background(), query, id).Scan(&u.ID, &u.Email, &u.FullName)

	if err != nil {
		return nil, err
	}
	return u, nil

}
