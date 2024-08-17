package author

import (
	"database/sql"
)

type Author struct {
	Id        int    `json:"id"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
}

type AuthorTable struct {
	db *sql.DB
}

func NewAuthorTable(db *sql.DB) AuthorTable {
	return AuthorTable{db: db}
}

func (t *AuthorTable) Get(id int) (Author, error) {
	var a Author
	if err := t.db.QueryRow("SELECT id, first_name, last_name FROM authors WHERE id = ?;", id).Scan(
		&a.Id,
		&a.Firstname,
		&a.Lastname,
	); err != nil {
		return Author{}, err
	}
	return a, nil
}

func (t *AuthorTable) GetAll() ([]Author, error) {
	rows, err := t.db.Query("SELECT id, first_name, last_name FROM authors;")
	if err != nil {
		return nil, err
	}

	authors := make([]Author, 0)
	for {
		if !rows.Next() {
			break
		}

		var a Author

		rows.Scan(&a.Id, &a.Firstname, &a.Lastname)

		authors = append(authors, a)
	}

	return authors, rows.Err()
}

func (t *AuthorTable) Insert(author Author) (Author, error) {
	row := t.db.QueryRow("INSERT INTO authors (first_name,last_name) VALUES (?,?) RETURNING id,first_name,last_name;", author.Firstname, author.Lastname)

	var a Author
	if err := row.Scan(&a.Id, &a.Firstname, &a.Lastname); err != nil {
		return Author{}, err
	}

	return a, nil
}
