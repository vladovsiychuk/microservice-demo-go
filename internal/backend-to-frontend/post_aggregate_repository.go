package backendtofrontend

type PostAggregateRepository struct {
}

type PostAggregateRepositoryI interface {
}

func NewPostAggregateRepository() *PostAggregateRepository {
	return &PostAggregateRepository{}
}
