# Unit Test Coverage Report

รายงานสรุปผลการทดสอบ (Unit Test Coverage) ของโปรเจกต์ CS367-G7-FoodDelivery

> การทดสอบทั้งหมดเป็น **Unit Test** โดยแต่ละ layer ทดสอบแบบ isolated ผ่านการ mock dependency (mockRepository, mockService, sqlmock)

---

## สรุป Coverage ระดับ Module

| Module | Coverage |
| :--- | :--- |
| **Menu Module** | 100.0% |
| **Middleware** | 100.0% |
| **JWT Package** | 100.0% |
| **Restaurant Module** | 99.0% |
| **Auth Module** | 96.6% |
| **Order Module** | 93.3% |

**Total Overall Coverage: 96.4%**

---

## รายละเอียด Unit Tests แยกตาม Layer

### 1. Handler Layer (HTTP Handler Unit Tests)

ทดสอบโดย mock Service — ตรวจสอบ HTTP status code และ response body

**Auth** (`internal/auth/auth_test.go`)

| Test Function | Test Cases | Coverage |
| :--- | :--- | :--- |
| `LoginHandler` | success, invalid JSON, user not found, wrong password | **100%** |

**Menu** (`internal/menu/handler_test.go`)

| Test Function | Test Cases | Coverage |
| :--- | :--- | :--- |
| `CreateMenu` | success (201), invalid restaurant ID, invalid JSON, service error (500) | **100%** |
| `GetMenu` | success (200), invalid restaurant ID, service error (500) | **100%** |

**Order** (`internal/order/handler_test.go`)

| Test Function | Test Cases | Coverage |
| :--- | :--- | :--- |
| `CreateOrder` | success, no user in context, invalid JSON | **76.5%** |
| `CancelOrder` | success, invalid ID, order not found, unauthorized | **100%** |
| `GetOrderByID` | success, invalid ID, order not found | **100%** |
| `UpdateOrderStatus` | success (restaurant/rider), invalid ID, no role, invalid JSON, not found, forbidden role, invalid transition | **95.5%** |
| `AssignRider` | success, invalid JSON, service error | **100%** |

**Restaurant** (`internal/restaurant/handler_test.go`)

| Test Function | Test Cases | Coverage |
| :--- | :--- | :--- |
| `CreateRestaurant` | success, bad JSON, missing field, no username, validation error, service error | **100%** |
| `GetRestaurants` | success, empty list returns `[]`, service error | **100%** |
| `GetRestaurantByID` | success, not found, invalid ID, service error | **100%** |
| `ConfirmOrder` | success, bad JSON, no username, order not found, service error | **100%** |

---

### 2. Service Layer (Business Logic Unit Tests)

ทดสอบโดย mock Repository — ตรวจสอบ business logic และ error handling

**Auth** (`internal/auth/auth_test.go`)

| Test Function | Test Cases |
| :--- | :--- |
| `Login` | success, user not found, wrong password, repository error |

**Menu** (`internal/menu/menu_test.go`)

| Test Function | Test Cases |
| :--- | :--- |
| `CreateMenu` | success, repository error |
| `GetMenu` | success, repository error |

**Order** (`internal/order/order_test.go`)

| Test Function | Test Cases |
| :--- | :--- |
| `CreateOrder` | success, food item not found, food item not available, repo error |
| `GetOrderByID` | success |
| `UpdateOrderStatus` | success (restaurant/rider), order not found, invalid role, invalid transition |
| `AssignRider` | success, repository error |

**Restaurant** (`internal/restaurant/restaurant_test.go`)

| Test Function | Test Cases |
| :--- | :--- |
| `CreateRestaurant` | success, invalid input (3 cases), repo error |
| `GetRestaurants` | returns list, empty list is not nil, repo error |
| `GetRestaurantByID` | success, not found, repo error |
| `ConfirmOrder` | success, repo error |

---

### 3. Repository Layer (Data Access Unit Tests)

ทดสอบโดย `go-sqlmock` — mock database driver โดยไม่ต่อ DB จริง

**Order** (`internal/order/order_test.go`)

| Test Function | Test Cases |
| :--- | :--- |
| `GetFoodItem` | success, not found (ErrNoRows) |
| `InsertOrderWithItems` | success (with transaction: begin, insert orders, insert items, commit) |
