# IMETRAX Backend API Documentation

Base URL: `http://localhost:8080/api/v1`

Auth header (protected routes):

```
Authorization: Bearer <token>
Content-Type: application/json
```

---

## 1) Produk & Alur Bisnis

- **Nama**: IMETRAX (Injection Mold Engineering TRAX)
- **Fungsi**: Digitalisasi produksi IME Formulatrix — scheduling part ke mesin CNC, pelacakan posisi & waktu, pembuatan Operation Plan + G-code, approval online (paperless), Gantt/PPIC scheduling, eksekusi operator.
- **Flow utama**: Engineer submit kebutuhan → **PEM** membuat Operation Plan + upload G-code → Approvals oleh **PEM, PPIC, QC, Engineering, Toolpather** (Manufacture Leader dapat diwakili role Engineering/Admin) → **PPIC** membuat jadwal mesin (PPIC schedule/Gantt) → **Operator** mengeksekusi (start/finish).
- **Status kunci**: Operation Plan `draft` → `pending_approval` → `approved`; G-code hanya di `draft`; assignment mesin `pending|in_progress|completed`.
- **Rate limit**: Auth 5 req/min; lainnya 100 req/min. Saat limit: `{"error":"Too many requests. Please try again later."}`

Peran ringkas:

- Admin: kelola user & mesin.
- PEM/Engineering: buat & submit Operation Plan + upload G-code.
- Approver (PEM, PPIC, QC, Engineering, Toolpather/Manufacture Leader): approve sampai lengkap.
- PPIC: jadwalkan mesin (PPIC schedule/Gantt).
- Operator: start/finish eksekusi.
- Guest: akses minimal.

---

## 2) Auth

### POST /auth/login _(public)_

- **Headers**: `Content-Type: application/json`
- **Body**:

  ```json
  { "user_id": "PI0824.2374", "password": "your_password" }
  ```

- **Response 200**:

  ```json
  {"message":"Login successful","data":{"token":"<jwt>","user":{...}}}
  ```

### POST /auth/register _(public – biasanya hanya dipakai admin awal/pengujian)_

- **Headers**: `Content-Type: application/json`
- **Roles valid**: `Admin`,`PPIC`,`Toolpather`,`PEM`,`QC`,`Engineering`,`Guest`,`Operator`
- **Body**:

  ```json
  {
    "user_id": "PI0824.0001",
    "username": "JOHN",
    "email": "john@company.com",
    "password": "password123",
    "role": "PPIC",
    "operator": "Operator Name"
  }
  ```

- **Response 201**: `{"message":"Registration successful","data":{...}}`

### GET /auth/profile _(protected)_

- **Headers**: `Authorization: Bearer <token>`
- **Body**: none
- **Response 200**: `{"data":{...user...}}`

### POST /auth/logout _(protected)_

- **Headers**: `Authorization: Bearer <token>`
- **Body**: none
- **Response 200**: `{"message":"Logout successful"}`

---

## 3) Machines

### GET /machines _(protected)_

- **Headers**: `Authorization: Bearer <token>`
- **Body**: none
- **Response**: `{"machines":[...],"count":N}`

### GET /machines/:id _(protected)_

- **Headers**: `Authorization: Bearer <token>`
- **Body**: none
- **Response**: `Machine`

### POST /admin/machines _(protected, Admin)_

- **Headers**: `Authorization: Bearer <token>`
- **Body**:

  ```json
  {
    "machine_code": "MC002",
    "machine_name": "CNC 2",
    "machine_type": "CNC",
    "location": "B",
    "status": "active"
  }
  ```

### PUT /admin/machines/:id _(protected, Admin)_

- **Headers**: `Authorization: Bearer <token>`
- **Body** (semua optional):

  ```json
  {
    "machine_name": "Updated",
    "machine_type": "5-axis",
    "location": "B",
    "status": "inactive"
  }
  ```

### DELETE /admin/machines/:id _(protected, Admin)_

- **Headers**: `Authorization: Bearer <token>`
- **Body**: none
- **Response 200**: `{"message":"Machine deleted successfully"}`

---

## 4) Job Orders

### GET /job-orders _(protected)_

- **Headers**: `Authorization: Bearer <token>`
- **Body**: none
- **Response**: `{"job_orders":[...],"count":N}`

### GET /job-orders/:id _(protected)_

- **Headers**: `Authorization: Bearer <token>`
- **Body**: none
- **Response**: `JobOrder` + `process_stages`

### GET /job-orders/machine/:machine*id *(protected)\_

- **Headers**: `Authorization: Bearer <token>`
- **Body**: none
- **Response**: list by machine

### POST /job-orders _(protected)_

- **Headers**: `Authorization: Bearer <token>`
- **Body**:

  ```json
  {
    "machine_id": 1,
    "njo": "NJO-2024-002",
    "project": "Project Beta",
    "item": "Part B",
    "note": "Optional",
    "deadline": "2024-12-31",
    "operator_id": 1
  }
  ```

### PUT /job-orders/:id _(protected)_

- **Headers**: `Authorization: Bearer <token>`
- **Body** (optional):

  ```json
  {
    "project": "Updated",
    "item": "Updated",
    "note": "Updated",
    "deadline": "2025-01-15",
    "operator_id": 2,
    "status": "completed"
  }
  ```

### DELETE /job-orders/:id _(protected)_

- **Headers**: `Authorization: Bearer <token>`
- **Body**: none
- **Response 200**: `{"message":"Job order deleted successfully"}`

---

## 5) Process Stages

### PUT /process-stages/:id _(protected)_

- **Headers**: `Authorization: Bearer <token>`
- **Body**:

  ```json
  {
    "start_time": "2024-01-01T09:00:00Z",
    "finish_time": "2024-01-01T17:00:00Z",
    "operator_id": 1,
    "notes": "Stage notes"
  }
  ```

- **Response 200**: `ProcessStage`

---

## 6) Operation Plans (approval & eksekusi)

Approver roles wajib: `PEM`, `PPIC`, `QC`, `Engineering`, `Toolpather` (Manufacture Leader dapat memakai Engineering/Admin).

### POST /operation-plans

```json
{
  "job_order_id": 1,
  "machine_id": 1,
  "part_quantity": 100,
  "description": "Plan desc"
}
```

Response 201: `{"message":"Operation plan created successfully","data":{...,"status":"draft"}}`

### GET /operation-plans

Query: `status` (`draft|pending_approval|approved`), `machine_id`
Response: list of plans

### GET /operation-plans/:id

Response: `{"data": OperationPlan}`

### DELETE /operation-plans/:id (hanya `draft`)

Response: `{"message":"Operation plan deleted successfully"}`

### POST /operation-plans/:id/submit

- Syarat: minimal 1 G-code terunggah
- Transisi: `draft` → `pending_approval`
  Response: `{"message":"Operation plan submitted for approval","data":...}`

### POST /operation-plans/:id/approve

- Hanya role: `PEM`,`PPIC`,`QC`,`Engineering`,`Toolpather`
- Tiap role approve sekali; lengkap → `approved` + email ke creator
  Response: `{"message":"Operation plan approved successfully","data":...}`

### GET /operation-plans/pending-approvals

Response: `{"data":[...plans...]}` (catatan: implementasi kini berbasis creator; perlu per-role enhancement)

### POST /operation-plans/:id/start (Operator)

- Syarat: status `approved`
  Response: `{"message":"Execution started successfully","data":...}`

### POST /operation-plans/:id/finish (Operator)

- Syarat: sudah start
  Response: `{"message":"Execution finished successfully","data":...}`

### GET /operation-plans/:id/approvals

Response: list `OperationPlanApproval`

### POST /operation-plans/:id/send-reminder (PEM)

- Kirim email reminder ke approver yang belum approve (butuh SMTP config)
  Response: counts sent/failed

---

## 7) G-Code Files

### POST /g-codes/upload

Form-data:

- `operation_plan_id` (int, required)
- `file` (.txt only, required)
- Syarat: plan `draft`
  Response 201: `{"message":"File uploaded successfully","data":{...}}`

### GET /g-codes/plan/:plan_id

Response: `{"data":[GCodeFile...]}`

### GET /g-codes/:id/download

Returns file attachment.

### DELETE /g-codes/:id

- Hanya saat plan `draft`
  Response: `{"message":"File deleted successfully"}`

---

## 8) Gantt Chart

### GET /gantt-chart

Query: `start_date`,`end_date` (YYYY-MM-DD), `priority` (`Low|Medium|Urgent|Top Urgent`), `status` (`pending|in_progress|completed`), `machine_id`, `group_by` (`priority|machine|`)
Response (ringkas):

```json
{
  "success": true,
  "data": {
    "sections": [{ "section_id":"priority_urgent","tasks":[...] }],
    "machines": [...],
    "summary": {...},
    "filters_applied": {...}
  }
}
```

---

## 9) PPIC Schedules (basis Gantt)

### GET /ppic-schedules

### GET /ppic-schedules/:id

### POST /ppic-schedules

```json
{
  "njo": "NJO-2024-001",
  "part_name": "Part A",
  "priority": "Urgent",
  "priority_alpha": "A",
  "material_status": "Ready",
  "start_date": "2024-01-01",
  "finish_date": "2024-01-10",
  "ppic_notes": "Notes",
  "machine_assignments": [
    {
      "machine_id": 1,
      "sequence": 1,
      "target_hours": 8.5,
      "scheduled_start": "2024-01-01T08:00:00Z",
      "scheduled_end": "2024-01-01T16:30:00Z"
    }
  ]
}
```

### PUT /ppic-schedules/:id

### DELETE /ppic-schedules/:id

### GET /ppic-schedules/machine/:machine_id

### POST /ppic-schedules/:id/machines

```json
{
  "machine_id": 2,
  "sequence": 2,
  "target_hours": 4.0,
  "scheduled_start": "2024-01-02T08:00:00Z",
  "scheduled_end": "2024-01-02T12:00:00Z"
}
```

### DELETE /ppic-schedules/:id/machines/:assignment_id

### PUT /ppic-schedules/:id/machines/:assignment_id/status

```json
{
  "status": "in_progress",
  "actual_start": "2024-01-01T08:15:00Z",
  "actual_end": null
}
```

Status: `pending|in_progress|completed`

---

## 10) Admin (role: Admin)

### GET /admin/users

Response: `{"users":[...],"total":N}`

### PUT /admin/users/:id

Request (opsional):

```json
{ "password": "new", "role": "Engineering", "is_active": true }
```

### DELETE /admin/users/:id

Response: `{"message":"User deleted successfully","username":"..."}` (tidak bisa self-delete)

---

## 11) Email

### GET /email/status

Response:

```json
{
  "success": true,
  "is_configured": true,
  "message": "Email service is configured and ready"
}
```

### POST /operation-plans/:id/send-reminder (lihat Operation Plans)

Env SMTP: `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM_ADDR`, `SMTP_FROM_NAME`, `FRONTEND_URL`

---

## 12) Error Format

- Umum: `{"error":"<message>"}`
- Validasi: `{"error":"Invalid request format","details":"..."}`
- Kode umum: 400, 401, 403, 404, 409, 500

---

## 13) Ringkasan Flow per Peran

- Engineer/PEM: buat Operation Plan → upload G-code → submit.
- Approver (PEM, PPIC, QC, Engineering, Toolpather/Manufacture Leader): approve hingga lengkap → status `approved` → email ke creator.
- PPIC: setelah approved, buat/update PPIC Schedule (mesin, prioritas, durasi) → terlihat di Gantt.
- Operator: start/finish eksekusi pada plan approved.
- Admin: kelola user & mesin.

---

## 14) Contoh Header

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

---

## 15) Catatan Implementasi

- G-code: hanya `.txt`; upload/delete hanya saat plan `draft`.
- Approval: 5 peran wajib; Manufacture Leader dapat menggunakan role Engineering/Admin bila perlu.
- Reminder email perlu SMTP aktif.
- Rate limiting aktif sesuai ketentuan di atas.
