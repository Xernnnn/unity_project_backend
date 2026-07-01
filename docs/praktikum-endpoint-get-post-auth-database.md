# Learning Session Backend: Praktikum Endpoint GET dan POST dengan Auth dan Database

Dokumen ini dipakai untuk workspace Learning Session Backend yang sudah dipoloskan sebagai bahan praktikum.

Di awal praktikum, endpoint berikut belum tersedia:

```txt
POST /learning-notes
GET  /learning-notes
```

Tugas peserta adalah membuat resource baru bernama `LearningNote` sampai endpoint tersebut bisa digunakan.

Resource ini sengaja sederhana agar fokus praktikum tetap jelas:

- Membuat endpoint `POST` untuk menyimpan data.
- Membuat endpoint `GET` untuk mengambil data.
- Memakai JWT auth yang sudah ada di project.
- Mengambil `user_id` dari auth middleware.
- Menyimpan dan membaca data menggunakan GORM dan PostgreSQL.
- Mengikuti pola folder project: `database`, `dto`, `contract`, `repository`, `service`, dan `controller`.

## 1. Gambaran Project

Project ini sudah menyediakan fondasi berikut:

- Server HTTP menggunakan Gin.
- Koneksi database PostgreSQL menggunakan GORM.
- Endpoint auth:

```txt
POST /auth/register
POST /auth/login
GET  /auth/profile
```

- Middleware auth yang membaca header:

```txt
Authorization: Bearer <token>
```

Saat praktikum dimulai, peserta tidak perlu membuat sistem login dari nol. Peserta cukup memakai auth yang sudah tersedia untuk mendapatkan token.

## 2. Alur Request

Alur endpoint yang akan dibuat:

```txt
Client / Postman / curl
    |
    v
Auth Middleware
    |
    v
Controller
    |
    v
Service
    |
    v
Repository
    |
    v
Database
```

Penjelasan singkat:

- `Auth Middleware` mengecek token dan menyimpan `user_id` ke Gin context.
- `Controller` membaca request, mengambil `user_id`, dan mengirim response JSON.
- `Service` berisi logic utama.
- `Repository` berisi query GORM.
- `Database model` menentukan struktur tabel.

## 3. Prasyarat Lokal

Pastikan komputer peserta sudah memiliki:

1. Go sesuai versi project.

```bash
go version
```

Project memakai Go `1.23.3` di `go.mod`. Versi Go yang lebih baru biasanya tetap bisa dipakai.

2. PostgreSQL.

```bash
psql --version
```

3. Tool untuk test HTTP.

Peserta boleh memakai salah satu:

- Postman
- Insomnia
- curl

4. OpenSSL untuk membuat key file.

```bash
openssl version
```

Jika OpenSSL tidak tersedia di Windows, peserta bisa memakai Git Bash, WSL, atau key file yang sudah disediakan instruktur.

## 4. Setup Workspace dari Awal

Masuk ke folder workspace yang diberikan instruktur:

```bash
cd learning-session-backend
```

Jika nama folder di komputer peserta berbeda, gunakan nama folder tersebut. Yang penting terminal berada di root project, yaitu folder yang berisi `go.mod`.

Install dependency Go:

```bash
go mod download
```

Jika perintah ini gagal karena internet atau proxy, minta bantuan instruktur. Dependency harus berhasil diunduh sebelum server bisa dijalankan.

## 5. Setup Database Praktikum

Buat database PostgreSQL khusus praktikum. Contoh nama database:

```txt
learning_session_backend_praktikum
```

Contoh memakai `psql`:

```bash
psql -U postgres
```

Lalu jalankan:

```sql
CREATE DATABASE learning_session_backend_praktikum;
```

Keluar dari `psql`:

```sql
\q
```

Catatan:

- Nama database boleh berbeda.
- Pastikan nama database di `.env` nanti sama dengan database yang dibuat.
- Praktikum ini hanya membutuhkan tabel `users` di awal, lalu peserta akan menambahkan tabel `learning_notes`.

## 6. Setup File .env

Di root project, buat atau edit file:

```txt
.env
```

Isi minimal yang dibutuhkan untuk praktikum:

```env
PORT=8080
IS_PRODUCTION=false
BASE_URL=http://localhost:8080

DB_USER=postgres
DB_PASS=ISI_PASSWORD_POSTGRES_ANDA
DB_NAME=learning_session_backend_praktikum
DB_HOST=127.0.0.1
DB_PORT=5432
DB_TIME_ZONE=Asia/Jakarta

PRIVATE_KEY=private_key.pem
PUBLIC_KEY=public_key.pem

ACCESS_TOKEN_LIFE_TIME=3600
REFRESH_TOKEN_LIFE_TIME=86400

ALLOW_ORIGIN=*
RATE_LIMIT_RPS=10
RATE_LIMIT_BURST=20
```

Yang wajib disesuaikan peserta:

- `DB_USER`
- `DB_PASS`
- `DB_NAME`
- `DB_HOST`
- `DB_PORT`

Catatan penting:

- Project membaca nama env `ACCESS_TOKEN_LIFE_TIME`, bukan `ACCESS_TOKEN_LIFETIME`.
- Project membaca nama env `REFRESH_TOKEN_LIFE_TIME`, bukan `REFRESH_TOKEN_LIFETIME`.
- File `.env` tidak perlu di-commit.

## 7. Setup JWT Key File

Project mewajibkan file berikut saat aplikasi start:

```txt
private_key.pem
public_key.pem
```

Jika kedua file tersebut sudah ada di root project, peserta boleh melewati langkah ini.

Jika belum ada, buat dengan OpenSSL:

```bash
openssl genrsa -out private_key.pem 2048
openssl rsa -in private_key.pem -pubout -out public_key.pem
```

Catatan:

- Pada versi praktikum ini, login masih memakai secret sementara di kode auth.
- Walaupun begitu, config tetap membaca file RSA key saat aplikasi start.
- Karena itu `private_key.pem` dan `public_key.pem` tetap harus tersedia.

## 8. Cek Server Awal

Sebelum menulis endpoint baru, jalankan migrasi awal:

```bash
go run main.go migrate
```

Jika berhasil, output akan menunjukkan migrasi dan seeding selesai.

Lalu jalankan server:

```bash
go run main.go
```

Server berjalan di:

```txt
http://localhost:8080
```

Jika port `8080` sudah dipakai, ubah `PORT` di `.env`, misalnya:

```env
PORT=8081
BASE_URL=http://localhost:8081
```

## 9. Cek Auth yang Sudah Ada

Saat seed dijalankan, project membuat admin default:

```txt
Email: admin@learningsession.local
Password: admin123
```

Login:

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@learningsession.local",
    "password": "admin123"
  }'
```

Response sukses berisi token:

```json
{
    "success": true,
    "message": "Login successful",
    "data": {
        "token": "TOKEN_ADA_DI_SINI",
        "user": {
            "id": 1,
            "name": "Admin",
            "email": "admin@learningsession.local",
            "role": "admin"
        }
    }
}
```

Simpan nilai `data.token`. Token ini akan dipakai untuk test endpoint `learning-notes`.

Cek profile:

```bash
curl -X GET http://localhost:8080/auth/profile \
  -H "Authorization: Bearer TOKEN_LOGIN"
```

Jika profile berhasil, setup awal sudah benar.

## 10. Spesifikasi Endpoint Praktikum

### 10.1 POST /learning-notes

Dipakai untuk membuat catatan belajar milik user yang sedang login.

Header:

```txt
Authorization: Bearer <token>
Content-Type: application/json
```

Request body:

```json
{
    "title": "Belajar GET dan POST",
    "content": "Hari ini saya belajar membuat endpoint dengan Gin dan GORM."
}
```

Response sukses:

```json
{
    "success": true,
    "message": "Learning note created successfully",
    "data": {
        "id": 1,
        "user_id": 1,
        "title": "Belajar GET dan POST",
        "content": "Hari ini saya belajar membuat endpoint dengan Gin dan GORM.",
        "created_at": "2026-06-04 10:00:00",
        "updated_at": "2026-06-04 10:00:00"
    }
}
```

### 10.2 GET /learning-notes

Dipakai untuk mengambil semua catatan belajar milik user yang sedang login.

Header:

```txt
Authorization: Bearer <token>
```

Response sukses:

```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "user_id": 1,
            "title": "Belajar GET dan POST",
            "content": "Hari ini saya belajar membuat endpoint dengan Gin dan GORM.",
            "created_at": "2026-06-04 10:00:00",
            "updated_at": "2026-06-04 10:00:00"
        }
    ]
}
```

## 11. Step 1: Buat Model Database

Buka file:

```txt
database/models.go
```

Tambahkan struct ini di bawah model `User`:

```go
type LearningNote struct {
    ID        int       `gorm:"column:id;primaryKey;autoIncrement;not null;<-:create"`
    UserID    int       `gorm:"column:user_id;not null;index"`
    Title     string    `gorm:"column:title;not null"`
    Content   string    `gorm:"column:content;type:text;not null"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

    User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}
```

Pastikan `database/models.go` memiliki import `time`. Pada workspace ini import tersebut sudah ada.

Penjelasan field:

- `ID`: primary key.
- `UserID`: pemilik catatan, diambil dari token login.
- `Title`: judul catatan.
- `Content`: isi catatan.
- `CreatedAt`: waktu data dibuat.
- `UpdatedAt`: waktu data terakhir diubah.
- `User`: relasi ke tabel `users`.

## 12. Step 2: Daftarkan Model ke Migration

Buka file:

```txt
database/migration.go
```

Cari:

```go
if err := db.AutoMigrate(
    &User{},
); err != nil {
```

Ubah menjadi:

```go
if err := db.AutoMigrate(
    &User{},
    &LearningNote{},
); err != nil {
```

Tujuannya agar GORM membuat tabel `learning_notes`.

## 13. Step 3: Buat DTO

Buat file baru:

```txt
dto/learning_note.go
```

Isi file:

```go
package dto

type CreateLearningNoteRequest struct {
    Title   string `json:"title" binding:"required,min=3"`
    Content string `json:"content" binding:"required,min=5"`
}

type LearningNoteListResponse struct {
    Success bool               `json:"success"`
    Data    []LearningNoteData `json:"data"`
}

type CreateLearningNoteResponse struct {
    Success bool             `json:"success"`
    Message string           `json:"message"`
    Data    LearningNoteData `json:"data"`
}

type LearningNoteData struct {
    ID        int    `json:"id"`
    UserID    int    `json:"user_id"`
    Title     string `json:"title"`
    Content   string `json:"content"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}
```

Penjelasan:

- `CreateLearningNoteRequest` membaca body `POST`.
- `LearningNoteListResponse` menjadi response `GET`.
- `CreateLearningNoteResponse` menjadi response `POST`.
- `LearningNoteData` menentukan bentuk JSON yang dikirim ke client.

## 14. Step 4: Tambahkan Contract Repository

Buka file:

```txt
contract/repository.go
```

Pada struct `Repository`, tambahkan field baru:

```go
LearningNote LearningNoteRepository
```

Hasilnya:

```go
type Repository struct {
    Auth         AuthRepository
    LearningNote LearningNoteRepository
}
```

Tambahkan interface baru di bawah `AuthRepository`:

```go
type LearningNoteRepository interface {
    Create(note *database.LearningNote) error
    FindByUserID(userID int) ([]database.LearningNote, error)
}
```

Penjelasan:

- `Create` dipakai untuk `POST /learning-notes`.
- `FindByUserID` dipakai untuk `GET /learning-notes`.
- Query memakai `userID` supaya user hanya melihat catatan miliknya sendiri.

## 15. Step 5: Buat Repository

Buat file baru:

```txt
repository/learning_note.go
```

Isi file:

```go
package repository

import (
    "github.com/RaFYWStud/LearningSessionBackend/contract"
    "github.com/RaFYWStud/LearningSessionBackend/database"
    "gorm.io/gorm"
)

type learningNoteRepository struct {
    db *gorm.DB
}

func ImplLearningNoteRepository(db *gorm.DB) contract.LearningNoteRepository {
    return &learningNoteRepository{db: db}
}

func (r *learningNoteRepository) Create(note *database.LearningNote) error {
    return r.db.Create(note).Error
}

func (r *learningNoteRepository) FindByUserID(userID int) ([]database.LearningNote, error) {
    var notes []database.LearningNote

    err := r.db.Where("user_id = ?", userID).
        Order("created_at DESC").
        Find(&notes).Error

    return notes, err
}
```

Penjelasan:

- `Create` melakukan insert ke tabel `learning_notes`.
- `FindByUserID` mengambil data milik user login.
- `Order("created_at DESC")` menampilkan catatan terbaru di atas.

## 16. Step 6: Daftarkan Repository

Buka file:

```txt
repository/repository.go
```

Ubah function `New` menjadi:

```go
func New(db *gorm.DB) *contract.Repository {
    return &contract.Repository{
        Auth:         ImplAuthRepository(db),
        LearningNote: ImplLearningNoteRepository(db),
    }
}
```

Jika langkah ini lupa, service tidak akan mendapat akses ke `LearningNoteRepository`.

## 17. Step 7: Tambahkan Contract Service

Buka file:

```txt
contract/service.go
```

Pada struct `Service`, tambahkan field baru:

```go
LearningNote LearningNoteService
```

Hasilnya:

```go
type Service struct {
    Auth         AuthService
    LearningNote LearningNoteService
}
```

Tambahkan interface baru:

```go
type LearningNoteService interface {
    CreateLearningNote(userID int, req dto.CreateLearningNoteRequest) (*dto.CreateLearningNoteResponse, error)
    GetMyLearningNotes(userID int) (*dto.LearningNoteListResponse, error)
}
```

Penjelasan:

- `CreateLearningNote` dipakai oleh endpoint POST.
- `GetMyLearningNotes` dipakai oleh endpoint GET.

## 18. Step 8: Buat Service

Buat file baru:

```txt
service/learning_note.go
```

Isi file:

```go
package service

import (
    "github.com/RaFYWStud/LearningSessionBackend/config/pkg/errs"
    "github.com/RaFYWStud/LearningSessionBackend/contract"
    "github.com/RaFYWStud/LearningSessionBackend/database"
    "github.com/RaFYWStud/LearningSessionBackend/dto"
)

type learningNoteService struct {
    learningNoteRepo contract.LearningNoteRepository
}

func ImplLearningNoteService(learningNoteRepo contract.LearningNoteRepository) contract.LearningNoteService {
    return &learningNoteService{learningNoteRepo: learningNoteRepo}
}

func (s *learningNoteService) CreateLearningNote(userID int, req dto.CreateLearningNoteRequest) (*dto.CreateLearningNoteResponse, error) {
    note := &database.LearningNote{
        UserID:  userID,
        Title:   req.Title,
        Content: req.Content,
    }

    if err := s.learningNoteRepo.Create(note); err != nil {
        return nil, errs.InternalServerError("failed to create learning note")
    }

    return &dto.CreateLearningNoteResponse{
        Success: true,
        Message: "Learning note created successfully",
        Data:    mapLearningNoteToDTO(*note),
    }, nil
}

func (s *learningNoteService) GetMyLearningNotes(userID int) (*dto.LearningNoteListResponse, error) {
    notes, err := s.learningNoteRepo.FindByUserID(userID)
    if err != nil {
        return nil, errs.InternalServerError("failed to fetch learning notes")
    }

    data := make([]dto.LearningNoteData, len(notes))
    for i, note := range notes {
        data[i] = mapLearningNoteToDTO(note)
    }

    return &dto.LearningNoteListResponse{
        Success: true,
        Data:    data,
    }, nil
}

func mapLearningNoteToDTO(note database.LearningNote) dto.LearningNoteData {
    return dto.LearningNoteData{
        ID:        note.ID,
        UserID:    note.UserID,
        Title:     note.Title,
        Content:   note.Content,
        CreatedAt: note.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt: note.UpdatedAt.Format("2006-01-02 15:04:05"),
    }
}
```

Penjelasan:

- Service menerima `userID` dari controller.
- Service membuat model `database.LearningNote`.
- Service memanggil repository untuk insert dan select.
- `mapLearningNoteToDTO` mengubah model database menjadi response JSON.

## 19. Step 9: Daftarkan Service

Buka file:

```txt
service/service.go
```

Ubah function `New` menjadi:

```go
func New(repo *contract.Repository) *contract.Service {
    return &contract.Service{
        Auth:         ImplAuthService(repo.Auth),
        LearningNote: ImplLearningNoteService(repo.LearningNote),
    }
}
```

Jika langkah ini lupa, controller tidak bisa memanggil `LearningNoteService`.

## 20. Step 10: Buat Controller

Buat file baru:

```txt
controller/learning_note.go
```

Isi file:

```go
package controller

import (
    "net/http"

    "github.com/RaFYWStud/LearningSessionBackend/config/middleware"
    "github.com/RaFYWStud/LearningSessionBackend/config/pkg/errs"
    "github.com/RaFYWStud/LearningSessionBackend/contract"
    "github.com/RaFYWStud/LearningSessionBackend/dto"
    "github.com/gin-gonic/gin"
)

type LearningNoteController struct {
    service contract.LearningNoteService
}

func (lc *LearningNoteController) GetPrefix() string {
    return "/learning-notes"
}

func (lc *LearningNoteController) InitService(service *contract.Service) {
    lc.service = service.LearningNote
}

func (lc *LearningNoteController) InitRoute(app *gin.RouterGroup) {
    auth := app.Group("")
    auth.Use(middleware.Auth())
    {
        auth.GET("", lc.getMyLearningNotes)
        auth.POST("", lc.createLearningNote)
    }
}

func (lc *LearningNoteController) getMyLearningNotes(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        HandlerError(ctx, errs.Unauthorized("user not authenticated"))
        return
    }

    response, err := lc.service.GetMyLearningNotes(userID.(int))
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, response)
}

func (lc *LearningNoteController) createLearningNote(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        HandlerError(ctx, errs.Unauthorized("user not authenticated"))
        return
    }

    var payload dto.CreateLearningNoteRequest
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid request payload"))
        return
    }

    response, err := lc.service.CreateLearningNote(userID.(int), payload)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusCreated, response)
}
```

Penjelasan:

- `GetPrefix` membuat prefix route `/learning-notes`.
- `middleware.Auth()` membuat semua route di controller ini wajib login.
- `ctx.Get("user_id")` mengambil user login dari token.
- `ShouldBindJSON` membaca body untuk endpoint POST.
- `http.StatusCreated` mengirim status `201 Created`.

## 21. Step 11: Daftarkan Controller

Buka file:

```txt
controller/controller.go
```

Tambahkan controller baru ke `allController`:

```go
allController := []Controller{
    &AuthController{},
    &LearningNoteController{},
    // Add your controller here
}
```

Jika langkah ini lupa, endpoint `/learning-notes` akan menghasilkan `404 Not Found`.

## 22. Step 12: Format Kode

Jalankan:

```bash
gofmt -w database/models.go database/migration.go dto/learning_note.go contract/repository.go contract/service.go repository/learning_note.go repository/repository.go service/learning_note.go service/service.go controller/learning_note.go controller/controller.go
```

Tujuannya agar format kode Go konsisten.

## 23. Step 13: Cek Build

Jalankan:

```bash
go test ./...
```

Jika berhasil, semua package akan muncul dengan status `no test files` atau `ok`.

Jika muncul error `undefined`, biasanya ada file yang belum dibuat atau belum didaftarkan.

## 24. Step 14: Jalankan Migration Ulang

Setelah model `LearningNote` didaftarkan, jalankan:

```bash
go run main.go migrate
```

Migration ini akan membuat tabel `learning_notes`.

Jika ingin mengulang database dari kosong:

```bash
go run main.go reset
```

Peringatan:

- `reset` menghapus tabel yang didaftarkan di migration.
- Gunakan hanya jika data praktikum boleh hilang.

## 25. Step 15: Jalankan Server

Jalankan:

```bash
go run main.go
```

Server:

```txt
http://localhost:8080
```

Biarkan terminal server tetap menyala. Buka terminal baru untuk menjalankan curl.

## 26. Step 16: Login untuk Mendapatkan Token

Login admin:

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@learningsession.local",
    "password": "admin123"
  }'
```

Ambil token dari:

```txt
data.token
```

Pada contoh berikutnya, ganti `TOKEN_LOGIN` dengan token tersebut.

## 27. Step 17: Test POST /learning-notes

```bash
curl -X POST http://localhost:8080/learning-notes \
  -H "Authorization: Bearer TOKEN_LOGIN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Belajar GET dan POST",
    "content": "Hari ini saya belajar membuat endpoint dengan Gin dan GORM."
  }'
```

Yang harus dicek:

- HTTP status `201 Created`.
- `success` bernilai `true`.
- `user_id` otomatis berasal dari token.
- Request body tidak mengirim `user_id`.

## 28. Step 18: Test GET /learning-notes

```bash
curl -X GET http://localhost:8080/learning-notes \
  -H "Authorization: Bearer TOKEN_LOGIN"
```

Yang harus dicek:

- HTTP status `200 OK`.
- Data yang dibuat lewat POST muncul.
- Urutan data terbaru berada di atas.

## 29. Step 19: Test Validasi Auth

Jalankan tanpa token:

```bash
curl -X GET http://localhost:8080/learning-notes
```

Response yang diharapkan:

```json
{
    "success": false,
    "error": "Authorization header required",
    "code": "MISSING_TOKEN"
}
```

Artinya middleware auth bekerja.

## 30. Step 20: Test Validasi Body

Kirim body yang tidak valid:

```bash
curl -X POST http://localhost:8080/learning-notes \
  -H "Authorization: Bearer TOKEN_LOGIN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "AB",
    "content": "ok"
  }'
```

Response yang diharapkan:

```json
{
    "status": 400,
    "error": "Bad Request",
    "message": "invalid request payload"
}
```

Artinya validasi DTO bekerja.

## 31. Step 21: Test Data Per User

Buat user baru:

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Peserta Dua",
    "email": "peserta2@example.com",
    "password": "password123",
    "password_confirmation": "password123"
  }'
```

Login sebagai user baru:

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "peserta2@example.com",
    "password": "password123"
  }'
```

Gunakan token user baru untuk:

```bash
curl -X GET http://localhost:8080/learning-notes \
  -H "Authorization: Bearer TOKEN_USER_BARU"
```

Yang harus dicek:

- Data admin tidak muncul untuk user baru.
- Jika user baru membuat catatan, catatan itu tidak muncul saat admin login.

Ini membuktikan filter `FindByUserID` bekerja.

## 32. Urutan Singkat Praktikum

Ikuti urutan ini saat mengerjakan:

1. Setup `.env`.
2. Pastikan `private_key.pem` dan `public_key.pem` tersedia.
3. Buat database PostgreSQL.
4. Jalankan `go mod download`.
5. Jalankan `go run main.go migrate`.
6. Jalankan `go run main.go`.
7. Test login auth.
8. Tambahkan model `LearningNote`.
9. Daftarkan model ke migration.
10. Buat DTO.
11. Tambahkan contract repository.
12. Buat repository.
13. Daftarkan repository.
14. Tambahkan contract service.
15. Buat service.
16. Daftarkan service.
17. Buat controller.
18. Daftarkan controller.
19. Jalankan `gofmt`.
20. Jalankan `go test ./...`.
21. Jalankan migration ulang.
22. Jalankan server.
23. Login untuk mendapatkan token.
24. Test `POST /learning-notes`.
25. Test `GET /learning-notes`.
26. Test error tanpa token.
27. Test error body tidak valid.
28. Test isolasi data antar user.

## 33. Checklist Peserta

- [ ] Saya berhasil menjalankan server sebelum membuat endpoint baru.
- [ ] Saya berhasil login dan mendapat token.
- [ ] Saya membuat model `LearningNote`.
- [ ] Saya mendaftarkan model ke migration.
- [ ] Saya membuat DTO request dan response.
- [ ] Saya membuat repository.
- [ ] Saya mendaftarkan repository.
- [ ] Saya membuat service.
- [ ] Saya mendaftarkan service.
- [ ] Saya membuat controller.
- [ ] Saya mendaftarkan controller.
- [ ] Saya menjalankan `gofmt`.
- [ ] Saya menjalankan `go test ./...`.
- [ ] Saya menjalankan migration ulang.
- [ ] Saya berhasil membuat data dengan POST.
- [ ] Saya berhasil mengambil data dengan GET.
- [ ] Saya memahami kenapa `user_id` tidak dikirim dari body request.
- [ ] Saya memastikan user hanya melihat catatan miliknya sendiri.

## 34. Kesalahan yang Sering Terjadi

### .env Tidak Terbaca

Gejala:

```txt
Environment variable is not set
```

Solusi:

- Pastikan file bernama `.env`.
- Pastikan file berada di root project.
- Pastikan nama env sesuai contoh, terutama `ACCESS_TOKEN_LIFE_TIME`.

### Database Tidak Bisa Connect

Gejala:

```txt
error connecting to SQL
```

Solusi:

- Pastikan PostgreSQL menyala.
- Pastikan `DB_USER`, `DB_PASS`, `DB_NAME`, `DB_HOST`, dan `DB_PORT` benar.
- Pastikan database sudah dibuat.

### Key File Tidak Ada

Gejala:

```txt
Error reading public key file
```

atau:

```txt
Error reading private key file
```

Solusi:

- Pastikan `private_key.pem` dan `public_key.pem` ada di root project.
- Pastikan `.env` menunjuk nama file yang benar.

### Lupa Header Authorization

Gejala:

```txt
401 Unauthorized
```

Solusi:

```txt
Authorization: Bearer TOKEN_LOGIN
```

### Lupa Header Content-Type

Gejala:

Body JSON tidak terbaca atau validasi gagal.

Solusi:

```txt
Content-Type: application/json
```

### Lupa Daftarkan Controller

Gejala:

```txt
404 Not Found
```

Solusi:

Tambahkan:

```go
&LearningNoteController{},
```

ke `allController`.

### Lupa Daftarkan Service

Gejala:

```txt
panic: runtime error: invalid memory address or nil pointer dereference
```

Solusi:

Tambahkan:

```go
LearningNote: ImplLearningNoteService(repo.LearningNote),
```

di `service/service.go`.

### Lupa Daftarkan Repository

Gejala:

Service tidak bisa memakai repository, atau field `repo.LearningNote` tidak ada.

Solusi:

Tambahkan:

```go
LearningNote: ImplLearningNoteRepository(db),
```

di `repository/repository.go`.

### Lupa Migration

Gejala:

```txt
relation "learning_notes" does not exist
```

Solusi:

Jalankan:

```bash
go run main.go migrate
```

### Salah Import

Gejala:

```txt
imported and not used
```

Solusi:

Jalankan:

```bash
gofmt -w nama_file.go
```

Lalu hapus import yang memang tidak dipakai.

## 35. Ringkasan Akhir

Endpoint yang dibuat:

```txt
POST /learning-notes
GET  /learning-notes
```

Keduanya memakai:

- Auth middleware.
- `user_id` dari JWT context.
- PostgreSQL melalui GORM.
- Pola `controller -> service -> repository`.

Setelah selesai, peserta seharusnya memahami alur lengkap membuat resource baru di backend Go tanpa membuat sistem auth dari nol.

