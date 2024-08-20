package tables

import "github.com/LucDeCaf/go-simple-blog/models"

type Table[M models.Model, PK any] interface {
	Get(PK) (M, error)
	GetAll() ([]M, error)
	Insert(M) (M, error)
	Delete(PK) (M, error)
}
