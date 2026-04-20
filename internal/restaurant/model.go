package restaurant

// Restaurant represents a restaurant owned by a user in the system.
type Restaurant struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	OwnerUsername string `json:"owner_username"`
}

// IsValid checks that all required fields are filled in.
// Used by the service layer before persisting the restaurant.
func (r *Restaurant) IsValid() bool {
	return r.Name != "" && r.Address != ""
}

// CreateRestaurantRequest is the expected JSON body for POST /restaurant.
// The owner is derived from the JWT claims in the handler, not the body.
type CreateRestaurantRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}
