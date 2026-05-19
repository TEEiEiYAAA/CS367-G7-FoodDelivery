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
users
├── id (PK)
├── username
├── password
└── role  [customer | restaurant_owner | rider]

restaurants
├── id (PK)
├── name
├── address
└── owner_username

food_items
├── id (PK)
├── restaurant_id (FK → restaurants)
├── name
├── price
└── is_available

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
├── food_item_id (FK → food_items)
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

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (รองรับทั้ง Windows, macOS, Linux)
- [Docker Hub](https://hub.docker.com/) account (สำหรับดึง image)

> ไม่จำเป็นต้องติดตั้ง Go หรือ MySQL เพิ่มเติม — ทุกอย่างรันผ่าน Docker

---

### Docker Image

API Image พร้อมใช้งานบน DockerHub แล้ว ไม่ต้อง build เอง

```
chonrathan/cs367-food-delivery:latest
```

---

### วิธีรันด้วย Docker Compose

**1. Clone repository**
```bash
git clone https://github.com/CS367-G7/CS367-G7-FoodDelivery.git
cd CS367-G7-FoodDelivery
```

**2. เปิด Docker Desktop** ให้พร้อมใช้งานก่อน

**3. Start ระบบ**
```bash
docker compose up
```

Docker จะดึง API image จาก DockerHub และ start MySQL อัตโนมัติ

รอจนเห็นข้อความ:
```
api-1  | Successfully connected to the database!
api-1  | Server running on :8080
```

**4. เรียกใช้งาน API ได้ที่**
```
http://localhost:8080
```

---

### หยุดระบบ

```bash
docker compose down
```

หากต้องการลบข้อมูลใน database ด้วย:
```bash
docker compose down -v
```

---

### วิธีรันด้วย Docker Compose แบบ Build เอง (ไม่ใช้ DockerHub Image)

หากต้องการ build image จาก source code เอง ให้แก้ [docker-compose.yml](../docker-compose.yml) โดยแทนที่ `image:` ด้วย `build:`:

```yaml
api:
  build:
    context: .
    dockerfile: Dockerfile
```

จากนั้นรัน:
```bash
docker compose up --build
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

## API Testing (Postman)

ไฟล์ Postman Collection และ Environment อยู่ใน folder `postman/`

```
postman/
├── CS367-G7-FoodDelivery.postman_collection.json
└── CS367-G7.postman_environment.json
```

### วิธี Import เข้า Postman

1. เปิด **Postman**
2. คลิก **Import** แล้วเลือกไฟล์ทั้งสองจาก folder `postman/`
3. เลือก Environment **CS367-G7** ที่มุมขวาบน

### วิธีรัน Collection Runner

1. คลิกขวาที่ Collection **CS367-G7-FoodDelivery**
2. เลือก **Run collection**
3. เลือก Environment **CS367-G7**
4. คลิก **Run CS367-G7-FoodDelivery**

> ระบบต้องรันอยู่ก่อน (`docker compose up`) และ API พร้อมที่ `http://localhost:8080`

### ผลการทดสอบ

| Tests | Passed | Failed |
|-------|--------|--------|
| 23 | 23 | 0 |

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
