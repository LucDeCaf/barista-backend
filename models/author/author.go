package author

type Author struct {
	Id        int    `json:"id"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
}

func NewAuthor(id int, firstName, lastName string) Author {
	return Author{
		Id:        id,
		Firstname: firstName,
		Lastname:  lastName,
	}
}
