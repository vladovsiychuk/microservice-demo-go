package backendtofrontend

type BffService struct {
}

type BffServiceI interface {
}

func NewService() *BffService {
	return &BffService{}
}
