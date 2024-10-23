package comment

type CommentService struct {
}

func NewService() *CommentService {
	return &CommentService{}
}

func (s *CommentService) GetAllComments() []string {
	return []string{"Comment 1", "Comment 2"}
}

func (s *CommentService) CreateComment(comment string) string {
	return "Comment created"
}
