package order

type Service interface {
	CreateOrder(username string, req CreateOrderRequest) (int64, int, error)
	CancelOrder(username string, orderID int) error
	GetOrderByID(orderID int) (*Order, []OrderItem, error)
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

// CancelOrder ส่งต่อ logic ไปยัง repository
func (s *service) CancelOrder(username string, orderID int) error {
	return s.repo.CancelOrder(username, orderID)
}
func (s *service) GetOrderByID(orderID int) (*Order, []OrderItem, error) {
	return s.repo.GetOrderByID(orderID)
}
func (s *service) UpdateOrderStatus() {}
func (s *service) AssignRider(orderID string, riderID int) error {
	return s.repo.AssignRider(orderID, riderID)
}