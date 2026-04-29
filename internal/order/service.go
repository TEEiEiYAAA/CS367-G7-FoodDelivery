package order

type Service interface {
	CreateOrder()
	CancelOrder()
	GetOrderByID()
	UpdateOrderStatus()
	AssignRider(orderID string, riderID int) error
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
func (s *service) AssignRider(orderID string, riderID int) error {
	return s.repo.AssignRider(orderID, riderID)
}