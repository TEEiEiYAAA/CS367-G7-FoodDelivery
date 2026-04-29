package order

type Service interface {
	CreateOrder(req CreateOrderRequest, customerUsername string) (Order, error)
	CancelOrder(orderID int, customerUsername string) error
	GetOrderByID(orderID int) (*Order, error)
	UpdateOrderStatus(orderID string, status string) error
	AssignRider(orderID string, riderID int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateOrder แปลง request เป็น domain objects แล้วให้ repo จัดการ
// การดึงราคาและคำนวณ subtotal / total_price ทำใน repo (ภายใน transaction เดียวกัน)
func (s *service) CreateOrder(req CreateOrderRequest, customerUsername string) (Order, error) {
	items := make([]OrderItem, len(req.Items))
	for i, r := range req.Items {
		items[i] = OrderItem{
			FoodItemID: r.FoodItemID,
			Quantity:   r.Quantity,
			// Subtotal จะถูก set โดย repo หลังดึงราคาจาก food_items
		}
	}

	order := Order{
		CustomerUsername: customerUsername,
		RestaurantID:     req.RestaurantID,
		DeliveryAddress:  req.DeliveryAddress,
		Status:           "pending",
	}

	return s.repo.CreateOrder(order, items)
}

func (s *service) CancelOrder(orderID int, customerUsername string) error {
	return s.repo.CancelOrder(orderID, customerUsername)
}

func (s *service) GetOrderByID(orderID int) (*Order, error) {
	return s.repo.GetOrderByID(orderID)
}

func (s *service) UpdateOrderStatus(orderID string, status string) error {
	return s.repo.UpdateOrderStatus(orderID, status)
}

func (s *service) AssignRider(orderID string, riderID int) error {
	return s.repo.AssignRider(orderID, riderID)
}