package users

import "github.com/LucDeCaf/go-simple-blog/db"

type User struct {
	Username             string `json:"username"`
	PasswordHashWithSalt string `json:"password_hash_with_salt"`
	Role                 Role   `json:"role"`
}

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func Get(username string) (User, error) {
	const query string = "SELECT username,password_hash_with_salt,role FROM users WHERE username=?;"

	row := db.DB.QueryRow(query, username)

	var u User
	err := row.Scan(
		&u.Username,
		&u.PasswordHashWithSalt,
		&u.Role,
	)

	return u, err
}

func GetAll() ([]User, error) {
	const query string = "SELECT username, password_hash_with_salt, role FROM users;"

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)
	for {
		if !rows.Next() {
			break
		}

		var u User

		if err := rows.Scan(
			&u.Username,
			&u.PasswordHashWithSalt,
			&u.Role,
		); err != nil {
			return users, err
		}

		users = append(users, u)
	}

	return users, rows.Err()
}

func Delete(username string) (User, error) {
	const query string = "DELETE FROM users WHERE id=? RETURNING %v;"

	var u User
	err := db.DB.QueryRow(query, username).Scan(
		&u.Username,
		&u.PasswordHashWithSalt,
		&u.Role,
	)

	return u, err
}

func Insert(user User) (User, error) {
	const query string = "INSERT INTO users (username,password_hash_with_salt,role) VALUES (?,?,?) RETURNING username,password_hash_with_salt,role;"

	var u User
	err := db.DB.QueryRow(query, user.Username, user.PasswordHashWithSalt, u.Role).Scan(
		&u.Username,
		&u.PasswordHashWithSalt,
		&u.Role,
	)

	return u, err
}
