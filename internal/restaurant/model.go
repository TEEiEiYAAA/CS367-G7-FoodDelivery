package restaurant

type Restaurant struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	OwnerUsername string `json:"owner_username"`
}

func (r *Restaurant) IsValid() bool {
	return r.Name != "" && r.Address != ""
}

// POST /restaurant
type CreateRestaurantRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

// PUT /restaurant/order/confirm
type ConfirmOrderRequest struct {
	OrderID int `json:"order_id" binding:"required"`
}
