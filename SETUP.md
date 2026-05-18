# CS367-G7 Food Delivery System

ระบบ Food Delivery API สำหรับจัดการร้านอาหาร เมนู และคำสั่งซื้อ พัฒนาด้วย Go (Gin Framework) และ MySQL

---

## System Overview

```
┌─────────────┐     HTTP      ┌──────────────────┐     SQL      ┌───────────┐
│   Client    │ ──────────▶  │   Go API Server  │ ──────────▶  │   MySQL   │
│ (REST API)  │              │  (Gin, Port 8080) │              │  (3306)   │
└─────────────┘              └──────────────────┘              └───────────┘
```

### Tech Stack

| Layer      | Technology               |
|------------|--------------------------|
| Language   | Go 1.25                  |
| Framework  | Gin v1.12                |
| Database   | MySQL 8.0                |
| Auth       | JWT (golang-jwt/jwt v5)  |
| Container  | Docker / Docker Compose  |

### Architecture

ระบบใช้ layered architecture แบ่งเป็น 3 ชั้น โดยมี Model เป็น shared struct ที่ทุกชั้นใช้ร่วมกัน

```
Handler  →  Service  →  Repository  →  MySQL
   ↑            ↑            ↑
         Model (shared structs)
```

- **Handler** — รับ HTTP request, validate input, ส่ง response
- **Service** — business logic
- **Repository** — query ฐานข้อมูล
- **Model** — โครงสร้างข้อมูล (struct) ที่ทุกชั้นใช้ร่วมกัน ไม่ใช่ชั้นแยกต่างหาก

### Database Schema

```
restaurants
├── id (PK)
├── name
├── address
└── owner_username

menus
├── id (PK)
├── restaurant_id (FK → restaurants)
├── name
├── price
└── stock

orders
├── id (PK)
├── customer_username
├── restaurant_id (FK → restaurants)
├── rider_id
├── status  [pending | confirmed | preparing | ready | delivering | delivered | cancelled]
├── total_price
├── delivery_address
├── created_at
└── customer_grace_period_end

order_items
├── id (PK)
├── order_id (FK → orders)
├── food_item_id (FK → menus)
├── quantity
└── subtotal
```

---

## API Endpoints

### Authentication

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/login` | - | รับ JWT token |

**Login Request**
```json
{
  "username": "your_username",
  "password": "your_password"
}
```

**Login Response**
```json
{
  "token": "<JWT_TOKEN>"
}
```

> Endpoint ที่ต้องการ Auth ให้ส่ง header: `Authorization: Bearer <JWT_TOKEN>`

---

### Restaurant

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/restaurant` | ✅ | สร้างร้านอาหาร |
| GET | `/restaurant` | - | ดูร้านอาหารทั้งหมด |
| GET | `/restaurant/:id` | - | ดูข้อมูลร้านอาหาร |
| PUT | `/restaurant/order/confirm` | ✅ | ยืนยันออเดอร์ (ฝั่งร้าน) |

---

### Menu

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/restaurant/:id/menu` | ✅ | เพิ่มเมนูในร้าน |
| GET | `/restaurant/:id/menu` | - | ดูเมนูของร้าน |

---

### Order

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/order` | ✅ | สร้างคำสั่งซื้อ |
| GET | `/order/:id` | - | ดูรายละเอียดออเดอร์ |
| PUT | `/order/cancel` | ✅ | ลูกค้ายกเลิกออเดอร์ |
| PUT | `/order/:id/status` | ✅ | อัปเดตสถานะออเดอร์ |

---

### Rider

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/order/:id/assign-rider` | ✅ | มอบหมายไรเดอร์ให้ออเดอร์ |

---

## Installation & Running

### Prerequisites

- [Docker](https://www.docker.com/) และ Docker Compose

### วิธีรันด้วย Docker Compose (แนะนำ)

```bash
# 1. Clone repository
git clone https://github.com/CS367-G7/CS367-G7-FoodDelivery.git
cd CS367-G7-FoodDelivery

# 2. Start services (API + MySQL)
docker compose up --build
```

API จะพร้อมใช้งานที่ `http://localhost:8080`

หยุด services:
```bash
docker compose down
```

ลบ volume ของ database ด้วย:
```bash
docker compose down -v
```

---

### วิธีรันโดยไม่ใช้ Docker (Manual)

**Prerequisites:** Go 1.21+, MySQL 8.0

**1. ตั้งค่า MySQL**

```sql
CREATE DATABASE food_delivery_db;
```

จากนั้น import schema:
```bash
mysql -u root -p food_delivery_db < docker/init.sql
```

**2. ตั้งค่า Environment Variables**

```bash
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=password
export DB_NAME=food_delivery_db
```

**3. รัน Server**

```bash
go mod download
go run ./cmd/server
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `127.0.0.1` | MySQL host |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `root` | MySQL username |
| `DB_PASSWORD` | `password` | MySQL password |
| `DB_NAME` | `food_delivery_db` | Database name |
| `GIN_MODE` | `debug` | Gin mode (`debug` / `release`) |

---

## Running Tests

```bash
go test ./...
```

รัน test พร้อม coverage:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```
