package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"asset-flow/internal/model"
)

var (
	ErrEmailExists       = errors.New("email already registered")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserNotPromotable = errors.New("user is not an active employee eligible for promotion")
)
type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}


func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	query := `
		INSERT INTO users (name, email, password, role, gender, department_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		u.Name, u.Email, u.Password, u.Role, u.Gender, u.DepartmentID, u.Status,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrEmailExists
		}
		return err
	}

	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, name, email, password, role, gender, department_id, status,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL`

	return r.scanOne(r.db.QueryRow(ctx, query, email))
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, name, email, password, role, gender, department_id, status,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL`

	return r.scanOne(r.db.QueryRow(ctx, query, id))
}

func (r *UserRepository) scanOne(row pgx.Row) (*model.User, error) {
	var u model.User
	err := row.Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.Gender,
		&u.DepartmentID, &u.Status, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) List(ctx context.Context, departmentID *int64) ([]model.User, error) {
	query := `
		SELECT id, name, email, password, role, gender, department_id, status,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE deleted_at IS NULL`
	args := []any{}

	if departmentID != nil {
		args = append(args, *departmentID)
		query += ` AND department_id = $1`
	}
	query += ` ORDER BY name ASC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(
			&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.Gender,
			&u.DepartmentID, &u.Status, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
		); err != nil {
			return nil, err
		}
		u.Password = ""
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *UserRepository) UpdateRole(ctx context.Context, id int64, newRole string) (*model.User, error) {
	query := `
		UPDATE users
		SET role = $2, updated_at = now()
		WHERE id = $1 AND role = $3 AND deleted_at IS NULL
		RETURNING id, name, email, password, role, gender, department_id, status,
		          created_at, updated_at, deleted_at`

	u, err := r.scanOne(r.db.QueryRow(ctx, query, id, newRole, model.RoleEmployee))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotPromotable
		}
		return nil, err
	}
	return u, nil
}