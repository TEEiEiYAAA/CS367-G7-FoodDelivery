package order

type Service interface {
	CreateOrder()
	CancelOrder()
	GetOrderByID()
	UpdateOrderStatus()
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateOrder()       {}
func (s *service) CancelOrder()       {}
func (s *service) GetOrderByID()      {}
func (s *service) UpdateOrderStatus() {}
