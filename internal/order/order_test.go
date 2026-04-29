package order

import (
	"errors"
	"testing"
)

// สร้าง mockRepository ขึ้นมาเพื่อจำลองพฤติกรรมของฐานข้อมูล
type mockRepository struct {
	Repository
	err error // เราจะใช้ตัวแปรนี้กำหนดว่าอยากให้ Repo คืนค่า error หรือไม่
}

// จำลองฟังก์ชัน AssignRider
func (m *mockRepository) AssignRider(orderID string, riderID int) error {
	return m.err
}

// ต้องประกาศฟังก์ชันอื่นๆ ให้ครบตาม Interface (แม้จะไม่ได้ใช้ในเทสนี้)
func (m *mockRepository) CreateOrder(username string, req CreateOrderRequest) (int64, int, error) {
	return 0, 0, m.err
}
func (m *mockRepository) CancelOrder()       {}
func (m *mockRepository) GetOrderByID()      {}
func (m *mockRepository) UpdateOrderStatus() {}

func TestAssignRider(t *testing.T) {
	// Case 1: มอบหมายไรเดอร์สำเร็จ (Happy Path)
	t.Run("Success - Should return nil when repo success", func(t *testing.T) {
		mockRepo := &mockRepository{err: nil} // จำลองว่า DB ทำงานปกติ
		service := NewService(mockRepo)

		err := service.AssignRider("1", 101)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})

	// Case 2: เกิด Error จาก Database (Bad Path)
	t.Run("Failure - Should return error when repo fails", func(t *testing.T) {
		mockRepo := &mockRepository{err: errors.New("database connection failed")} // จำลอง DB พัง
		service := NewService(mockRepo)

		err := service.AssignRider("1", 101)

		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}