package menu

import (
	"errors"
	"testing"
)

type mockRepo struct {
	createMenuFn func(restaurantID int, req CreateMenuRequest) (Menu, error)
	getMenuFn    func(restaurantID int) ([]Menu, error)
}

func (m *mockRepo) CreateMenu(restaurantID int, req CreateMenuRequest) (Menu, error) {
	return m.createMenuFn(restaurantID, req)
}

func (m *mockRepo) GetMenu(restaurantID int) ([]Menu, error) {
	return m.getMenuFn(restaurantID)
}

func TestGetMenuByRestaurant(t *testing.T) {
	// TODO: Setup mock repository and test service
}

func TestCreateMenu(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		want := Menu{ID: 1, RestaurantID: 10, Name: "Pad Thai", Price: 80}
		repo := &mockRepo{
			createMenuFn: func(restaurantID int, req CreateMenuRequest) (Menu, error) {
				return want, nil
			},
		}
		svc := NewService(repo)
		got, err := svc.CreateMenu(10, CreateMenuRequest{Name: "Pad Thai", Price: 80})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got != want {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("repository_error", func(t *testing.T) {
		repo := &mockRepo{
			createMenuFn: func(restaurantID int, req CreateMenuRequest) (Menu, error) {
				return Menu{}, errors.New("db error")
			},
		}
		svc := NewService(repo)
		_, err := svc.CreateMenu(10, CreateMenuRequest{Name: "Pad Thai", Price: 80})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
