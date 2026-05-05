package menu

type Service interface {
	CreateMenu(restaurantID int, req CreateMenuRequest) (Menu, error)
	GetMenu()
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateMenu(restaurantID int, req CreateMenuRequest) (Menu, error) {
	return s.repo.CreateMenu(restaurantID, req)
}
func (s *service) GetMenu() {}
