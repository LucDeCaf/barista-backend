package user

import "database/sql"

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Salt         string `json:"salt"`
	Role         string `json:"role"`
}

type UserTable struct {
	db *sql.DB
}

func NewUserTable(db *sql.DB) UserTable {
	return UserTable{db: db}
}

func (t UserTable) Get(username string) (User, error) {
	row := t.db.QueryRow("SELECT username, password_hash, salt, role FROM users WHERE username=?;", username)

	var u User

	if err := row.Scan(&u.Username, &u.PasswordHash, &u.Salt, &u.Role); err != nil {
		return User{}, err
	}

	return u, nil
}

func (t UserTable) GetAll() ([]User, error) {
	rows, err := t.db.Query("SELECT username, password_hash, salt, role FROM users;")
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
			&u.PasswordHash,
			&u.Salt,
			&u.Role,
		)

		users = append(users, u)
	}

	return users, rows.Err()
}
