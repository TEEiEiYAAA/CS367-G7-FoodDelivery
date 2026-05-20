package menu

type Menu struct {
	ID           int     `json:"id"`
	RestaurantID int     `json:"restaurant_id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
}

type CreateMenuRequest struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
}
