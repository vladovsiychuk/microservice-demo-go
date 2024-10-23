package post

type PostService struct {
}

func NewService() *PostService {
	return &PostService{}
}

func (s *PostService) GetAllPosts() []string {
	return []string{"Post 1", "Post 2"}
}

func (s *PostService) CreatePost(req PostRequest) (*Post, error) {
	post, err := CreatePost(req)
	if err != nil {
		return nil, err
	}

	return post, nil
}
