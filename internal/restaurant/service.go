package restaurant

type Service interface {
	CreateRestaurant()
	GetRestaurants()
	GetRestaurantByID()
	ConfirmOrder()
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateRestaurant()  {}
func (s *service) GetRestaurants()    {}
func (s *service) GetRestaurantByID() {}
func (s *service) ConfirmOrder()      {}
