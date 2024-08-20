package user

import (
	"database/sql"
	"fmt"
)

// TODO: This is error-prone to Scan from, possibly create func `ScanDefault(*sql.Row, *Blog) error`
const PubliclyReturned string = "username,password_hash_with_salt,role"

type User struct {
	Username             string `json:"username"`
	PasswordHashWithSalt string `json:"password_hash_with_salt"`
	Role                 Role   `json:"role"`
}

type UserTable struct {
	db *sql.DB
}

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func NewUser(username, passwordHash string, role Role) *User {
	return &User{
		Username:             username,
		PasswordHashWithSalt: passwordHash,
		Role:                 role,
	}
}

func (u *User) Fields() []any {
	return []any{&u.Username, &u.PasswordHashWithSalt, &u.Role}
}

func (u *User) Values() []any {
	return []any{u.Username, u.PasswordHashWithSalt, u.Role}
}

func NewUserTable(db *sql.DB) UserTable {
	return UserTable{db: db}
}

func (t UserTable) Get(username string) (User, error) {
	row := t.db.QueryRow("SELECT username, password_hash_with_salt, role FROM users WHERE username=?;", username)

	var u User

	if err := row.Scan(&u.Username, &u.PasswordHashWithSalt, &u.Role); err != nil {
		return User{}, err
	}

	return u, nil
}

func (t UserTable) GetAll() ([]User, error) {
	rows, err := t.db.Query("SELECT username, password_hash_with_salt, role FROM users;")
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)
	for {
		if !rows.Next() {
			break
		}

		var u User

		rows.Scan(
			&u.Username,
			&u.PasswordHashWithSalt,
			&u.Role,
		)

		users = append(users, u)
	}

	return users, rows.Err()
}

func (t UserTable) Insert(user *User) (*User, error) {
	query := fmt.Sprintf("INSERT INTO users (username,password_hash_with_salt,role) VALUES (?,?,?) RETURNING %v;", PubliclyReturned)

	var u User

	row := t.db.QueryRow(query, user.Values()...)
	if err := row.Scan(u.Fields()...); err != nil {
		return nil, err
	}

	return &u, nil
}
