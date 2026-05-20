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

**Total Overall Coverage: 84.8%**

---

## รายละเอียด Unit Tests แยกตาม Layer

### 1. Handler Layer (HTTP Handler Unit Tests)

ทดสอบโดย mock Service — ตรวจสอบ HTTP status code และ response body

**Auth** (`internal/auth/auth_test.go`)

| Test Function | Test Cases |
| :--- | :--- |
| `LoginHandler` | success, invalid JSON, user not found, wrong password |

**Menu** (`internal/menu/handler_test.go`)

| Test Function | Test Cases |
| :--- | :--- |
| `CreateMenu` | success (201), invalid restaurant ID, invalid JSON, service error (500) |
| `GetMenu` | success (200), invalid restaurant ID, service error (500) |

**Order** (`internal/order/order_test.go`, `internal/order/handler_test.go`, `internal/order/cancel_order_test.go`)

| Test Function | Test Cases |
| :--- | :--- |
| `CreateOrder` | success, no user in context, invalid JSON |
| `CancelOrder` | success, unauthorized, invalid JSON, invalid token claims, order not found, forbidden, order cannot be cancelled (unprocessable), internal error |
| `GetOrderByID` | success, invalid ID, order not found |
| `UpdateOrderStatus` | success (restaurant/rider), invalid ID, no role, invalid JSON, not found, forbidden role, invalid transition |
| `AssignRider` | success, invalid JSON, service error |

**Restaurant** (`internal/restaurant/handler_test.go`)

| Test Function | Test Cases |
| :--- | :--- |
| `CreateRestaurant` | success, bad JSON, missing field, no username, validation error, service error |
| `GetRestaurants` | success, empty list returns `[]`, service error |
| `GetRestaurantByID` | success, not found, invalid ID, service error |
| `ConfirmOrder` | success, bad JSON, no username, order not found, service error |

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

**Order** (`internal/order/order_test.go`, `internal/order/cancel_order_test.go`)

| Test Function | Test Cases |
| :--- | :--- |
| `CreateOrder` | success, food item not found, food item not available, repo error |
| `CancelOrder` | success, order not found, forbidden, status not pending, grace period expired, update error |
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

### 3. Middleware & Package Unit Tests

**Middleware** (`internal/middleware/middleware_test.go`)

| Test Function | Test Cases |
| :--- | :--- |
| `AuthMiddleware` | valid token, no header, invalid token |
| `RequireRole` | correct role, no role in context, wrong role |

**JWT Package** (`pkg/jwt/jwt_test.go`)

| Test Function | Test Cases |
| :--- | :--- |
| `GenerateToken` | success |
| `ValidateToken` | success (valid token), invalid token, empty token |

---

### 4. Repository Layer (Data Access Unit Tests)

ทดสอบโดย `go-sqlmock` — mock database driver โดยไม่ต่อ DB จริง

**Auth** (`internal/auth/repository_test.go`)

| Test Function | Test Cases |
| :--- | :--- |
| `GetUserByUsername` | success, not found (ErrNoRows), DB error |

**Menu** (`internal/menu/repository_test.go`)

| Test Function | Test Cases |
| :--- | :--- |
| `CreateMenu` | success, exec error, last insert ID error |
| `GetMenu` | success, query error, scan error, rows error |

**Order** (`internal/order/order_test.go`, `internal/order/cancel_order_test.go`)

| Test Function | Test Cases |
| :--- | :--- |
| `GetFoodItem` | success, not found (ErrNoRows) |
| `InsertOrderWithItems` | success (with transaction: begin, insert orders, insert items, commit) |
| `GetOrder` | success, order not found (ErrNoRows), scan error |
| `SetOrderStatus` | success, database error |

**Restaurant** (`internal/restaurant/repository_test.go`)

| Test Function | Test Cases |
| :--- | :--- |
| `CreateRestaurant` | success, insert error, last insert ID error |
| `GetRestaurants` | success, empty is not nil, query error, scan error, rows error |
| `GetRestaurantByID` | success, not found (ErrNoRows) |
| `ConfirmOrder` | success, order not found, exec error |
