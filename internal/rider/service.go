package rider

type Service interface {
	AssignRider()
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) AssignRider() {}
