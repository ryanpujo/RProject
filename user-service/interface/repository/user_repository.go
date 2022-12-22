package repository

import (
	"context"
	"database/sql"
	"time"
	"user-service/models"
	"user-service/usecases/repository"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) Create(user *models.UserPayload) (id int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := `insert into users (first_name, last_name, username, email, password, created_at, updated_at)
								values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err = ur.db.QueryRowContext(ctx, statement,
		user.Fname,
		user.Lname,
		user.Username,
		user.Email,
		user.Password,
		time.Now(),
		time.Now(),
	).Scan(&id)

	return
}

func (ur *userRepository) FindById(id int64) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `select id, first_name, last_name, email, username, password, created_at, updated_at from users where id = $1`
	row := ur.db.QueryRowContext(ctx, stmt, id)
	var result models.User
	err := row.Scan(
		&result.Id,
		&result.Fname,
		&result.Lname,
		&result.Email,
		&result.Username,
		&result.Password,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (ur *userRepository) FindUsers() (users []*models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `select id, first_name, last_name, password, email, username, created_at, updated_at from users order by first_name`
	var rows *sql.Rows
	rows, err = ur.db.QueryContext(ctx, stmt)
	if err != nil {
		return
	}
	for rows.Next() {
		var result models.User
		err = rows.Scan(
			&result.Id,
			&result.Fname,
			&result.Lname,
			&result.Password,
			&result.Email,
			&result.Username,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return
		}
		users = append(users, &result)
	}
	return
}

func (ur *userRepository) DeleteById(id int64) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `delete from users where id=$1`

	_, err = ur.db.ExecContext(ctx, stmt, id)
	return
}

func (ur *userRepository) Update(user *models.UserPayload) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update users set
		first_name = $1,
		last_name = $2,
		username = $3,
		email = $4,
		password = $5,
		updated_at = $6
		where id = $7
	`
	_, err = ur.db.ExecContext(ctx, stmt,
		user.Fname,
		user.Lname,
		user.Username,
		user.Email,
		user.Password,
		time.Now(),
		user.Id,
	)
	return
}

func (ur *userRepository) FindByUsername(username string) (user *models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `select id, first_name, last_name, password, email, username, created_at, updated_at from users where username = $1`
	row := ur.db.QueryRowContext(ctx, stmt, username)
	var result models.User
	err = row.Scan(
		&result.Id,
		&result.Fname,
		&result.Lname,
		&result.Password,
		&result.Email,
		&result.Username,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	user = &result
	return
}
