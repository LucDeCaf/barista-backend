package author

import (
	"database/sql"
	"fmt"

	"github.com/LucDeCaf/go-simple-blog/models"
)

type Author struct {
	Id        int    `json:"id"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
}

type AuthorTable struct {
	db *sql.DB
}

func (*Author) FieldNames() string {
	return "id,first_name,last_name"
}

func FieldNames() string {
	return (*Author).FieldNames(nil)
}

func (a *Author) Fields() []any {
	return []any{&a.Id, &a.Firstname, &a.Lastname}
}

func (a *Author) Values() []any {
	return models.ValuesFromFields(a)
}

func NewAuthorTable(db *sql.DB) AuthorTable {
	return AuthorTable{db: db}
}

func (t AuthorTable) Get(id int) (*Author, error) {
	query := fmt.Sprintf("SELECT %v FROM authors WHERE id = ?;", FieldNames())

	var a Author

	row := t.db.QueryRow(query, id)
	if err := row.Scan(a.Fields()...); err != nil {
		return nil, err
	}

	return &a, nil
}

func (t AuthorTable) GetAll() ([]*Author, error) {
	query := fmt.Sprintf("SELECT %v FROM authors;", FieldNames())

	rows, err := t.db.Query(query)
	if err != nil {
		return nil, err
	}

	authors := make([]*Author, 0)
	for {
		if !rows.Next() {
			break
		}

		var a Author

		rows.Scan(&a.Id, &a.Firstname, &a.Lastname)

		authors = append(authors, &a)
	}

	return authors, rows.Err()
}

func (t AuthorTable) Insert(author *Author) (*Author, error) {
	query := fmt.Sprintf("INSERT INTO authors (first_name,last_name) VALUES (?,?) RETURNING %v;", FieldNames())

	var a Author

	row := t.db.QueryRow(query, author.Firstname, author.Lastname)
	if err := row.Scan(a.Fields()...); err != nil {
		return nil, err
	}

	return &a, nil
}

func (t AuthorTable) Delete(id int) (*Author, error) {
	query := fmt.Sprintf("DELETE FROM authors WHERE id=? RETURNING %v;", FieldNames())

	var a Author

	row := t.db.QueryRow(query, id)
	if err := row.Scan(&a.Id, &a.Firstname, &a.Lastname); err != nil {
		return nil, err
	}

	return &a, nil
}
