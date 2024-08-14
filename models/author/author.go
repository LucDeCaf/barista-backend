package author

type Author struct {
	Id        int    `json:"id"`
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
}

func NewAuthor(id int, firstName, lastName string) Author {
	return Author{
		Id:        id,
		Firstname: firstName,
		Lastname:  lastName,
	}
}
