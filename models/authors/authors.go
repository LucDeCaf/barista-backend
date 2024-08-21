package authors

import "github.com/LucDeCaf/go-simple-blog/db"

type Author struct {
	Id        int    `json:"id"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
}

func Get(id int) (Author, error) {
	const query string = "SELECT id,first_name,last_name FROM authors WHERE id=?;"

	var a Author
	err := db.DB.QueryRow(query, id).Scan(
		&a.Id,
		&a.Firstname,
		&a.Lastname,
	)

	return a, err
}

func GetAll() ([]Author, error) {
	const query string = "SELECT id,first_name,last_name FROM authors;"

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	authors := make([]Author, 0)
	for {
		if !rows.Next() {
			break
		}

		var a Author

		if err := rows.Scan(
			&a.Id,
			&a.Firstname,
			&a.Lastname,
		); err != nil {
			return authors, err
		}

		authors = append(authors, a)
	}

	return authors, rows.Err()
}

func Insert(author Author) (Author, error) {
	const query string = "INSERT INTO authors (first_name,last_name) VALUES (?,?) RETURNING id,first_name,last_name;"

	var a Author
	err := db.DB.QueryRow(query, author.Firstname, author.Lastname).Scan(
		&a.Id,
		&a.Firstname,
		&a.Lastname,
	)

	return a, err
}

func Update(author Author) (Author, error) {
	const query string = "INSERT INTO authors (first_name,last_name) VALUES (?,?) WHERE id=? RETURNING id,first_name,last_name;"

	var a Author
	err := db.DB.QueryRow(query, author.Firstname, author.Lastname).Scan(
		&a.Id,
		&a.Firstname,
		&a.Lastname,
	)

	return a, err
}

func Delete(id int) (Author, error) {
	const query string = "DELETE FROM authors WHERE id=? RETURNING id,first_name,last_name;"

	var a Author
	err := db.DB.QueryRow(query, id).Scan(
		&a.Id,
		&a.Firstname,
		&a.Lastname,
	)

	return a, err
}
