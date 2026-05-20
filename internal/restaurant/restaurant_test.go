package restaurant

import (
	"database/sql"
	"errors"
	"testing"
)

// mock repository
type mockRepo struct {
	createFn  func(Restaurant) (Restaurant, error)
	listFn    func() ([]Restaurant, error)
	getByIDFn func(int) (*Restaurant, error)
	confirmFn func(int, string) error
}

func (m *mockRepo) CreateRestaurant(r Restaurant) (Restaurant, error) { return m.createFn(r) }
func (m *mockRepo) GetRestaurants() ([]Restaurant, error)             { return m.listFn() }
func (m *mockRepo) GetRestaurantByID(id int) (*Restaurant, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(id)
	}
	return nil, nil
}
func (m *mockRepo) ConfirmOrder(orderID int, ownerUsername string) error {
	if m.confirmFn != nil {
		return m.confirmFn(orderID, ownerUsername)
	}
	return nil
}

// Restaurant.IsValid
func TestRestaurant_IsValid(t *testing.T) {
	tests := []struct {
		name string
		r    Restaurant
		want bool
	}{
		{"valid", Restaurant{Name: "KFC", Address: "Bangkok"}, true},
		{"missing name", Restaurant{Name: "", Address: "Bangkok"}, false},
		{"missing address", Restaurant{Name: "KFC", Address: ""}, false},
		{"both missing", Restaurant{}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.r.IsValid(); got != tc.want {
				t.Errorf("IsValid() = %v, want %v (input: %+v)", got, tc.want, tc.r)
			}
		})
	}
}

// Service.CreateRestaurant
func TestCreateRestaurant_Success(t *testing.T) {
	repo := &mockRepo{
		createFn: func(r Restaurant) (Restaurant, error) {
			r.ID = 42
			return r, nil
		},
	}
	svc := NewService(repo)

	input := Restaurant{Name: "KFC", Address: "Bangkok", OwnerUsername: "alice"}
	got, err := svc.CreateRestaurant(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 42 {
		t.Errorf("got ID %d, want 42", got.ID)
	}
	if got.Name != "KFC" || got.Address != "Bangkok" || got.OwnerUsername != "alice" {
		t.Errorf("unexpected restaurant: %+v", got)
	}
}

func TestCreateRestaurant_InvalidRejectsWithoutCallingRepo(t *testing.T) {
	cases := []struct {
		name  string
		input Restaurant
	}{
		{"missing name", Restaurant{Name: "", Address: "Bangkok"}},
		{"missing address", Restaurant{Name: "KFC", Address: ""}},
		{"both missing", Restaurant{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockRepo{
				createFn: func(r Restaurant) (Restaurant, error) {
					t.Fatalf("repo should not be called for invalid input: %+v", r)
					return Restaurant{}, nil
				},
			}
			svc := NewService(repo)

			_, err := svc.CreateRestaurant(tc.input)
			if !errors.Is(err, ErrInvalidRestaurant) {
				t.Errorf("got %v, want ErrInvalidRestaurant", err)
			}
		})
	}
}

func TestCreateRestaurant_RepoErrorIsPropagated(t *testing.T) {
	dbErr := errors.New("database is down")
	repo := &mockRepo{
		createFn: func(r Restaurant) (Restaurant, error) {
			return Restaurant{}, dbErr
		},
	}
	svc := NewService(repo)

	_, err := svc.CreateRestaurant(Restaurant{Name: "KFC", Address: "Bangkok"})
	if !errors.Is(err, dbErr) {
		t.Errorf("got %v, want %v", err, dbErr)
	}
}

// Service.GetRestaurants
func TestGetRestaurants_ReturnsList(t *testing.T) {
	want := []Restaurant{
		{ID: 1, Name: "KFC", Address: "Bangkok", OwnerUsername: "alice"},
		{ID: 2, Name: "MK", Address: "Nonthaburi", OwnerUsername: "bob"},
	}
	repo := &mockRepo{
		listFn: func() ([]Restaurant, error) {
			return want, nil
		},
	}
	svc := NewService(repo)

	got, err := svc.GetRestaurants()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %+v want %+v", i, got[i], want[i])
		}
	}
}

func TestGetRestaurants_EmptyListIsNotNil(t *testing.T) {
	repo := &mockRepo{
		listFn: func() ([]Restaurant, error) {
			return []Restaurant{}, nil
		},
	}
	svc := NewService(repo)

	got, err := svc.GetRestaurants()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Error("got nil slice, expected empty (non-nil) slice")
	}
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %d items", len(got))
	}
}

func TestGetRestaurants_RepoErrorIsPropagated(t *testing.T) {
	dbErr := errors.New("database is down")
	repo := &mockRepo{
		listFn: func() ([]Restaurant, error) {
			return nil, dbErr
		},
	}
	svc := NewService(repo)

	_, err := svc.GetRestaurants()
	if !errors.Is(err, dbErr) {
		t.Errorf("got %v, want %v", err, dbErr)
	}
}

// Service.GetRestaurantByID
func TestGetRestaurantByID_Success(t *testing.T) {
	want := &Restaurant{ID: 1, Name: "KFC", Address: "Bangkok", OwnerUsername: "alice"}
	repo := &mockRepo{
		getByIDFn: func(id int) (*Restaurant, error) { return want, nil },
	}
	svc := NewService(repo)

	got, err := svc.GetRestaurantByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *got != *want {
		t.Errorf("got %+v, want %+v", *got, *want)
	}
}

func TestGetRestaurantByID_NotFound(t *testing.T) {
	repo := &mockRepo{
		getByIDFn: func(id int) (*Restaurant, error) { return nil, sql.ErrNoRows },
	}
	svc := NewService(repo)

	_, err := svc.GetRestaurantByID(99)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("got %v, want ErrNotFound", err)
	}
}

func TestGetRestaurantByID_RepoErrorIsPropagated(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mockRepo{
		getByIDFn: func(id int) (*Restaurant, error) { return nil, dbErr },
	}
	svc := NewService(repo)

	_, err := svc.GetRestaurantByID(1)
	if !errors.Is(err, dbErr) {
		t.Errorf("got %v, want %v", err, dbErr)
	}
}

// Service.ConfirmOrder
func TestConfirmOrder_Success(t *testing.T) {
	repo := &mockRepo{
		confirmFn: func(orderID int, ownerUsername string) error { return nil },
	}
	svc := NewService(repo)

	if err := svc.ConfirmOrder(1, "alice"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestConfirmOrder_RepoErrorIsPropagated(t *testing.T) {
	repo := &mockRepo{
		confirmFn: func(orderID int, ownerUsername string) error { return ErrOrderNotFound },
	}
	svc := NewService(repo)

	err := svc.ConfirmOrder(99, "alice")
	if !errors.Is(err, ErrOrderNotFound) {
		t.Errorf("got %v, want ErrOrderNotFound", err)
	}
}
