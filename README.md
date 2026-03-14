# CS367-G7-FoodDelivery
### รายชื่อสมาชิก
<table border="1" cellpadding="8" cellspacing="1" width="100%">
  <thead>
    <tr>
      <th align="left">ชื่อ-นามสกุล</th>
      <th align="right">รหัสนักศึกษา</th>
      <th align="left">ชื่อ Account GitHub</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>นางสาวธวัลหทัย เทียมทอง</td>
      <td align="right">6609650111</td>
      <td><a href="https://github.com/tha15thai">tha15thai</a></td>
    </tr>
    <tr>
      <td>นางสาวกนกพร พรรณปัญญา</td>
      <td align="right">6609650152</td>
      <td><a href="https://github.com/Kanokporn-6609650152">NgerNgaEiiei</a></td>
    </tr>
    <tr>
      <td>นายกุลเศรษฐ์ เนตรเพชร</td>
      <td align="right">6609650194</td>
      <td><a href="https://github.com/JJJOBJY">JJJOBJY</a></td>
    </tr>
    <tr>
      <td>นางสาวชลธาร ศิลปาจารย์</td>
      <td align="right">6609650269</td>
      <td><a href="https://github.com/Chonr">Chonr</a></td>
    </tr>
    <tr>
      <td>นายธีรัตม์ ศรีสุโข</td>
      <td align="right">6609650442</td>
      <td><a href="https://github.com/TEEiEiYAAA">TEEiEiYAAA</a></td>
    </tr>
  </tbody>
</table>

---
## Feature
มีทั้งหมด 12 Feature
### 🏪 Restaurant
POST /restaurant (สร้างร้านอาหาร)

GET /restaurant (ดูร้านทั้งหมด)

GET /restaurant/{id} (ดูข้อมูลของร้านอาหาร)

PUT /restaurant/order/confirm (ยืนยันออเดอร์)

### 🍽 Menu
POST /restaurant/{id}/menu (เพิ่มเมนู)

GET /restaurant/{id}/menu (ดูเมนู)

PUT /restaurant/menu/stock (อัปเดตจำนวนสินค้า)

PUT/restaurant/{id}/menu/status (สถานะเมนู) 

### 🧾 Order
POST /order (สร้างคำสั่งซื้อ) 

PUT /order/cancel (ลูกค้ายกเลิกออเดอร์) 

GET /order/{id} (ดูรายละเอียดออเดอร์) 

PUT / order / {id} / status (อัปเดตสถานะออเดอร์ เช่น รับออเดอร์ กำลังทำ ทำเสร็จ กำลังจัดส่ง )

### 🛵 Rider
POST /order/{id}/assign-rider (มอบหมายไรเดอร์) 

---
### Responsibility

นางสาวธวัลหทัย เทียมทอง
- 
-

## นางสาวกนกพร พรรณปัญญา
- POST /restaurant (สร้างร้านอาหาร)
- GET /restaurant (ดูร้านทั้งหมด)
  
นายกุลเศรษฐ์ เนตรเพชร
-
-

นางสาวชลธาร ศิลปาจารย์
-
-

นายธีรัตม์ ศรีสุโข
-
-
