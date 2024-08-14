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

func (at AuthorTable) Get(id int) (Author, error) {
	var a Author
	if err := at.db.QueryRow("SELECT id, first_name, last_name FROM authors WHERE id = ?;", id).Scan(
		&a.Id,
		&a.Firstname,
		&a.Lastname,
	); err != nil {
		return Author{}, err
	}
	return a, nil
}

func (at AuthorTable) GetAll() ([]Author, error) {
	// Select rows
	rows, err := at.db.Query("SELECT id, first_name, last_name FROM authors;")
	if err != nil {
		return nil, err
	}

	// Read authors from rows into slice
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
