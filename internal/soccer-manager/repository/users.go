package repository

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -destination=mock/mock_user.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/repository UserRepository
type UserRepository interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error)
	DeleteUser(ctx context.Context, id int64) error
	GetUserByID(ctx context.Context, id int64) (GetUserByIDRow, error)
	GetUserByUsername(ctx context.Context, username string) (GetUserByUsernameRow, error)
	GetUserHashByID(ctx context.Context, id int64) (string, error)
	GetAuth(ctx context.Context, username string) (GetAuthRow, error)
	ListUsersCursor(ctx context.Context, arg ListUsersCursorParams) ([]ListUsersCursorRow, error)
	UpdateUserHash(ctx context.Context, arg UpdateUserHashParams) error
	CheckUserExists(ctx context.Context, id int64) (bool, error)
}

type pgUserRepo struct {
	db            *pgxpool.Pool
	snowflakeNode *snowflake.Node
}

func NewUserRepo(db *pgxpool.Pool, snowflakeNode *snowflake.Node) *pgUserRepo {
	return &pgUserRepo{
		db:            db,
		snowflakeNode: snowflakeNode,
	}
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, username, role, hash) VALUES ($1, $2, $3, $4) RETURNING id, username, role
`

type CreateUserParams struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Hash     string `json:"hash"`
}

type CreateUserRow struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (r *pgUserRepo) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := r.db.QueryRow(ctx, createUser,
		r.snowflakeNode.Generate().Int64(),
		arg.Username,
		arg.Role,
		arg.Hash,
	)
	var i CreateUserRow
	err := row.Scan(&i.ID, &i.Username, &i.Role)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1
`

func (r *pgUserRepo) DeleteUser(ctx context.Context, id int64) error {
	res, err := r.db.Exec(ctx, deleteUser, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, username, role FROM users WHERE id = $1
`

type GetUserByIDRow struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (r *pgUserRepo) GetUserByID(ctx context.Context, id int64) (GetUserByIDRow, error) {
	row := r.db.QueryRow(ctx, getUserByID, id)
	var i GetUserByIDRow
	err := row.Scan(&i.ID, &i.Username, &i.Role)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, role FROM users WHERE username = $1
`

type GetUserByUsernameRow struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (r *pgUserRepo) GetUserByUsername(
	ctx context.Context,
	username string,
) (GetUserByUsernameRow, error) {
	row := r.db.QueryRow(ctx, getUserByUsername, username)
	var i GetUserByUsernameRow
	err := row.Scan(&i.ID, &i.Username, &i.Role)
	return i, err
}

const getUserHashByID = `-- name: GetUserHashByID :one
SELECT hash FROM users WHERE id = $1 LIMIT 1
`

func (r *pgUserRepo) GetUserHashByID(ctx context.Context, id int64) (string, error) {
	row := r.db.QueryRow(ctx, getUserHashByID, id)
	var hash string
	err := row.Scan(&hash)
	return hash, err
}

const getAuth = `-- name: GetAuth :one
SELECT id, hash, role FROM users WHERE username = $1 LIMIT 1
`

type GetAuthRow struct {
	ID   int64  `json:"id"`
	Hash string `json:"hash"`
	Role string `json:"role"`
}

func (r *pgUserRepo) GetAuth(ctx context.Context, username string) (GetAuthRow, error) {
	row := r.db.QueryRow(ctx, getAuth, username)
	var i GetAuthRow
	err := row.Scan(&i.ID, &i.Hash, &i.Role)
	return i, err
}

const listUsersWithIDGreaterThan = `-- name: ListUsersWithIDGreaterThan :many
SELECT id, username FROM users WHERE id > $1 ORDER BY id LIMIT $2
`

const listUsersCursor = `-- name: ListUsersCursor :many
SELECT id, username, role FROM users WHERE id > $1 ORDER BY id LIMIT $2
`

type ListUsersCursorParams struct {
	ID    int64 `json:"id"`
	Limit int32 `json:"limit"`
}

type ListUsersCursorRow struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (r *pgUserRepo) ListUsersCursor(
	ctx context.Context,
	arg ListUsersCursorParams,
) ([]ListUsersCursorRow, error) {
	rows, err := r.db.Query(ctx, listUsersCursor, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUsersCursorRow
	for rows.Next() {
		var i ListUsersCursorRow
		if err := rows.Scan(&i.ID, &i.Username, &i.Role); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserHash = `-- name: UpdateUserHash :exec
UPDATE users SET hash = $2 WHERE id = $1
`

type UpdateUserHashParams struct {
	ID   int64  `json:"id"`
	Hash string `json:"hash"`
}

func (r *pgUserRepo) UpdateUserHash(ctx context.Context, arg UpdateUserHashParams) error {
	res, err := r.db.Exec(ctx, updateUserHash, arg.ID, arg.Hash)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

const checkUserExists = `-- name: CheckUserExists :one
SELECT EXISTS (SELECT 1 FROM users WHERE id = $1) AS user_exists
`

func (r *pgUserRepo) CheckUserExists(ctx context.Context, id int64) (bool, error) {
	row := r.db.QueryRow(ctx, checkUserExists, id)
	var user_exists bool
	err := row.Scan(&user_exists)
	return user_exists, err
}
