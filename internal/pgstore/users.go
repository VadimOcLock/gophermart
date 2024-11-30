package pgstore

import (
	"context"

	"github.com/VadimOcLock/gophermart/internal/entity"

	"github.com/VadimOcLock/gophermart/internal/service/authservice"
)

const createUser = `
INSERT INTO users (login, password_hash)
VALUES ($1, $2)
RETURNING id;
`

func (q *Queries) CreateUser(ctx context.Context, params authservice.CreateUserParams) (uint64, error) {
	row := q.db.QueryRow(ctx, createUser, params.Login, params.PasswordHash)
	var id uint64
	err := row.Scan(&id)

	return id, err
}

const userExistsByLogin = `
SELECT EXISTS(SELECT 1 FROM users WHERE login = $1);
`

func (q *Queries) UserExistsByLogin(ctx context.Context, login string) (bool, error) {
	row := q.db.QueryRow(ctx, userExistsByLogin, login)
	var exists bool
	err := row.Scan(&exists)

	return exists, err
}

const findUserByLogin = `
select id, login, password_hash, created_at 
FROM users
WHERE login = $1
limit 1
;
`

func (q *Queries) FindUserByLogin(ctx context.Context, login string) (entity.User, error) {
	row := q.db.QueryRow(ctx, findUserByLogin, login)
	var user entity.User
	err := row.Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	return user, err
}
