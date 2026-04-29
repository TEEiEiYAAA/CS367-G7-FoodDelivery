package order

type Service interface {
	CreateOrder(username string, req CreateOrderRequest) (int64, int, error)
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

// CreateOrder ส่งต่อ logic ไปยัง repository
func (s *service) CreateOrder(username string, req CreateOrderRequest) (int64, int, error) {
	return s.repo.CreateOrder(username, req)
}

func (s *service) CancelOrder()       {}
func (s *service) GetOrderByID()      {}
func (s *service) UpdateOrderStatus() {}
func (s *service) AssignRider(orderID string, riderID int) error {
	return s.repo.AssignRider(orderID, riderID)
}