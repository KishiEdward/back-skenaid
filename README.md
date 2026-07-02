# Skenaid Backend API

<div align="center">
  <img width="300" height="301" alt="Institut Teknologi dan Bisnis Bina Sarana Global" src="https://github.com/user-attachments/assets/1e84f66a-135b-4cf2-b07a-b2a9098ce119" width="200"/>
</div>

<div align="center">
Institut Teknologi dan Bisnis Bina Sarana Global <br>
FAKULTAS TEKNOLOGI INFORMASI & KOMUNIKASI <br>
https://global.ac.id/
</div>

## Project UAS
- Nim : 1123150045
- Nama : Dzidan Rafi Habibie
- Mata Kuliah : Pemrograman Mobile Lanjutan
- Kelas : TI-SE 23 M

## Deskripsi Project
Project ini adalah backend API untuk sistem Skenaid yang digunakan untuk mendukung autentikasi pengguna, manajemen produk, keranjang belanja, serta proses pemesanan (order) secara aman. Backend ini juga terintegrasi dengan aplikasi Skewallet melalui deep link untuk proses pembayaran. Aplikasi ini dibangun menggunakan Go, MySQL, dan Firebase Authentication.

## Link proyek lain yang terintegrasi
- **[Backend skenaid (ecommerce)](https://github.com/KishiEdward/back-skenaid)**
- **[Frontend skenaid (ecommerce)](https://github.com/KishiEdward/front-skenaid)**
- **[Backend skewallet (ewallet)](https://github.com/KishiEdward/back-skewallet)**
- **[Frontend skenaid (ewallet)](https://github.com/KishiEdward/front-skewallet)**

## Demo Video
Lihat demo aplikasi dan alur fitur yang tersedia dalam video berikut.

**[Watch Full Demo on YouTube]()**

Alternative link: **[Google Drive Demo]()**

## Fitur Utama
- Autentikasi pengguna dengan Firebase Authentication
- Registrasi dan login pengguna
- Manajemen data produk (CRUD)
- Keranjang belanja pengguna
- Manajemen pesanan (order) dan riwayat transaksi
- Deep link integration dengan aplikasi Skewallet untuk proses pembayaran
- Middleware autentikasi untuk melindungi endpoint
- Seeding data produk untuk kebutuhan development

## Teknologi yang Digunakan
- **[Go](https://go.dev/)** - Bahasa pemrograman backend
- **[MySQL](https://www.mysql.com/)** - Database relasional
- **[Firebase](https://firebase.google.com/)** - Authentication
- **[JWT](https://jwt.io/)** - Token autentikasi
- **[godotenv](https://github.com/joho/godotenv)** - Manajemen environment variable

## Persyaratan Sistem
Pastikan perangkat Anda sudah memiliki:
- Go (versi terbaru yang kompatibel dengan modul ini)
- MySQL Server
- Firebase project dengan service account
- Git
- Postman (opsional untuk testing API)

## Cara Menjalankan Project

### 1. Clone Repository
```bash
git clone https://github.com/KishiEdward/skenaid-backend.git
cd skenaid
```

### 2. Install Dependency
```bash
go mod tidy
```

### 3. Siapkan Environment
Salin file `.env.example` menjadi `.env`, lalu sesuaikan konfigurasinya, contohnya:
```bash
cp .env.example .env
```
```env
APP_PORT=8080
APP_ENV=development
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=skenaid
JWT_SECRET=1234567890abcdef1234567890abcdef
JWT_EXPIRE_HOURS=24
FIREBASE_CREDENTIALS_PATH=./firebase-service-account.json
GOOGLE_APPLICATION_CREDENTIALS=./firebase-service-account.json
```

### 4. Siapkan Firebase
- Buat project Firebase
- Aktifkan Authentication
- Download file service account JSON
- Simpan sebagai `firebase-service-account.json` di root project

### 5. Siapkan Database MySQL
Pastikan MySQL sudah berjalan dan database yang digunakan sudah tersedia sebelum menjalankan server.

### 6. Jalankan Server
```bash
go run main.go
```

Server akan berjalan di:
```bash
http://localhost:8080
```

### 7. Seed Data Produk (Opsional)
Untuk mengisi data produk contoh ke database:
```bash
go run ./seed
```

## Struktur Project
```bash
skenaid/
├── config/                          # Konfigurasi database dan Firebase
│   ├── database.go
│   └── firebase.go
├── handlers/                        # Handler HTTP untuk endpoint API
│   ├── auth_handler.go
│   ├── cart_handler.go
│   ├── order_handler.go
│   ├── product_handler.go
│   └── user_handler.go
├── middleware/                      # Middleware autentikasi
│   └── auth_middleware.go
├── models/                          # Struktur data model
│   ├── cart.go
│   ├── order.go
│   ├── product.go
│   └── user.go
├── repositories/                    # Layer akses data ke database
│   ├── product_repository.go
│   └── user_repository.go
├── routes/                          # Routing API
│   └── router.go
├── seed/                            # Script seeding data produk contoh
│   └── seed.go
├── services/                        # Logika bisnis aplikasi
│   ├── auth_service.go
│   └── product_service.go
├── main.go                          # Entry point aplikasi
├── go.mod                           # Modul Go
├── go.sum                           # Checksum dependency Go
├── .env                             # Konfigurasi environment
├── .gitignore
└── firebase-service-account.json    # Kredensial Firebase Admin SDK
```

## Dokumentasi API
Base URL:
```bash
http://localhost:8080/v1
```

### Authentication
- `POST /v1/auth/register` - Registrasi pengguna baru
- `POST /v1/auth/login` - Login pengguna dan menghasilkan token
- `POST /v1/auth/verify-token` - Verifikasi Firebase ID Token dan menghasilkan JWT backend

### Users
- `GET /v1/users/profile` - Ambil data profil pengguna (butuh token)
- `PUT /v1/users/profile` - Update data profil pengguna

### Products
- `GET /v1/products` - Ambil daftar produk
- `GET /v1/products/:id` - Ambil detail produk
- `POST /v1/products` - Tambah produk
- `PUT /v1/products/:id` - Update produk
- `DELETE /v1/products/:id` - Hapus produk

### Cart
- `GET /v1/cart` - Ambil isi keranjang pengguna
- `POST /v1/cart` - Tambah produk ke keranjang
- `PUT /v1/cart/:id` - Update quantity item keranjang
- `DELETE /v1/cart/:id` - Hapus satu item keranjang

### Orders
- `POST /v1/orders/checkout` - Checkout keranjang menjadi pesanan
- `GET /v1/orders` - Ambil riwayat pesanan pengguna
- `GET /v1/orders/:id` - Ambil detail pesanan

## Integrasi Deep Link
Backend ini mendukung alur pembayaran lintas aplikasi bersama Skewallet:
- Aplikasi Skenaid memanggil skema `skewallet://pay` untuk membuka aplikasi Skewallet dan memproses pembayaran.
- Setelah transaksi selesai, Skewallet mengarahkan callback kembali ke Skenaid melalui skema `skenaid://payment-callback`.

## Lisensi
Project ini dilisensikan di bawah MIT License.

## Ucapan Terima Kasih
- [Firebase](https://firebase.google.com/)
- [MySQL](https://www.mysql.com/)
- [Go](https://go.dev/)

---
<div align="center">
  <p>© 2026 Skenaid Backend API. All rights reserved.</p>
</div>
