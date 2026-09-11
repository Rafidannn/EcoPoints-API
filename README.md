# EcoPoints Golang REST API

RESTful API service untuk aplikasi **EcoPoints** (Bank Sampah Digital), dibangun menggunakan **Golang**, **Gin Framework**, **GORM**, dan terintegrasi dengan database MySQL.

---

## 🚀 Fitur Utama

- 🔐 **Authentication & Authorization**: JWT (JSON Web Token) dengan Role-Based Access Control (Admin, Petugas, Nasabah)
- ♻️ **Waste Types API**: CRUD data jenis sampah, kategori, dan nilai poin per kg
- 📍 **Drop Points API**: CRUD data lokasi bank sampah / drop point
- 🎁 **Rewards API**: CRUD data reward & penukaran poin
- 📄 **Swagger Documentation**: Dokumentasi interaktif OpenAPI / Swagger UI
- 🐳 **Docker Ready**: Dukungan multi-stage build Docker image yang ringan (~15MB)

---

## 🛠️ Tech Stack

- **Language**: Go 1.24+
- **Web Framework**: [Gin Web Framework](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/) (MySQL Driver)
- **Authentication**: `golang-jwt/jwt/v5` + `bcrypt`
- **Docs**: [swaggo/swag](https://github.com/swaggo/swag)

---

## 📋 Prerequisites

- Go `1.22+` / `1.24+`
- MySQL / MariaDB Server

---

## ⚙️ Setup & Menjalankan Lokal

### 1. Clone & Masuk Direktori
```bash
git clone <URL_REPO_ANDA>
cd ecopoints-go-api
```

### 2. Salin Konfigurasi Environment
```bash
cp .env.example .env
```
Sesuaikan konfigurasi koneksi database dan JWT di file `.env`:
```env
APP_NAME=EcoPoints-Go-API
APP_ENV=development
APP_PORT=8090

DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=ecopoints
DB_USERNAME=root
DB_PASSWORD=

JWT_SECRET=your_jwt_secret_key_here
JWT_EXPIRATION_HOURS=72
```

### 3. Download Dependencies
```bash
go mod download
```

### 4. Jalankan Aplikasi
```bash
go run cmd/api/main.go
```
Aplikasi akan berjalan di: `http://localhost:8090`

---

## 📖 Dokumentasi Swagger UI

Setelah aplikasi berjalan, buka browser dan akses:
👉 **[http://localhost:8090/swagger/index.html](http://localhost:8090/swagger/index.html)**

Untuk me-regenerate dokumentasi Swagger setelah mengubah anotasi API:
```bash
swag init -g cmd/api/main.go
```

---

## 🚢 Deployment

### Opsi 1: Docker
```bash
# Build image
docker build -t ecopoints-go-api .

# Run container
docker run -d -p 8090:8090 --env-file .env --name ecopoints-api ecopoints-go-api
```

### Opsi 2: Systemd Linux Service (VPS)
```bash
# Build binary
go build -o ecopoints-api cmd/api/main.go
```
Jalankan binary dengan systemd service atau supervisor di server Anda.
