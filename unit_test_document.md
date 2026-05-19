# Unit Test Coverage Report

รายงานสรุปผลการทดสอบ (Unit Test Coverage) ของโปรเจกต์ CS367-G7-FoodDelivery

## สรุป Coverage ระดับ API (API Handlers)

**1. Auth API (`internal/auth`)**
* `LoginHandler`: **100.0%**

**2. Menu API (`internal/menu`)**
* `CreateMenu`: **100.0%**
* `GetMenu`: **100.0%**

**3. Order API (`internal/order`)**
* `CreateOrder`: **76.5%**
* `CancelOrder`: **100.0%**
* `GetOrderByID`: **100.0%**
* `UpdateOrderStatus`: **95.5%**
* `AssignRider`: **100.0%**

**4. Restaurant API (`internal/restaurant`)**
* `CreateRestaurant`: **100.0%**
* `GetRestaurants`: **100.0%**
* `GetRestaurantByID` (เรียกร้านอาหารตาม ID): **100.0%**
* `ConfirmOrder` (ยืนยันออเดอร์โดยร้าน): **100.0%**

---

## สรุปภาพรวม Coverage ของแต่ละ Module (รวม Service, Repository)

| Module | Coverage |
| :--- | :--- |
| **Menu Module** | 100.0% |
| **Middleware** | 100.0% |
| **JWT Package** | 100.0% |
| **Restaurant Module**| 99.0% |
| **Auth Module** | 96.6% |
| **Order Module** | 93.3% |

**Total Overall Coverage : 84.8%**
