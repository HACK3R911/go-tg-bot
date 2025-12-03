package repository

type AuthRepo interface {
	AuthorizeRepo(userID int64)
	IsAuthorizedRepo(userID int64) bool
}

type Repository struct {
	AuthRepo
}

func NewRepository() *Repository {
	return &Repository{
		AuthRepo: NewAuthDB(),
	}
}
