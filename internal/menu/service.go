package menu

type Service interface {
	CreateMenu()
	GetMenu()
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateMenu() {}
func (s *service) GetMenu()    {}
