# Learning Session Backend

Backend sederhana untuk praktikum pembuatan endpoint menggunakan Go, Gin, GORM, PostgreSQL, dan JWT auth.

## Endpoint

Auth:

```txt
POST /auth/register
POST /auth/login
GET  /auth/profile
```

## Setup

Install dependencies:

```bash
go mod download
```

Pastikan `.env` sudah berisi konfigurasi database dan server.

Jalankan migrasi:

```bash
go run main.go migrate
```

Jalankan server:

```bash
go run main.go
```

Default server:

```txt
http://localhost:8080
```

## Default Admin

Saat seed dijalankan, sistem membuat admin default:

```txt
Email: admin@learningsession.local
Password: admin123
```

## Contoh Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@learningsession.local",
    "password": "admin123"
  }'
```

