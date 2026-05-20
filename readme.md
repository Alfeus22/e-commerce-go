# 🛒 E-Commerce API (Go Backend)

Deskripsi singkat tentang projekmu di sini. Contoh: RESTful API e-commerce yang dibangun menggunakan Go dan MongoDB, dirancang dengan arsitektur bersih untuk menangani manajemen produk, autentikasi multi-role (Customer, Seller, Courier), hingga sistem riwayat pesanan dinamis.

---

## 🚀 Fitur Utama

- **Autentikasi & Otorisasi:** Registrasi dan login menggunakan JWT (JSON Web Token) dengan pembagian role (*Seller, Customer, admin*).
- **Manajemen Produk (Seller):** CRUD produk, edit produk sendiri, proteksi akses, dan fitur unggah gambar produk.
- **Pencarian Fleksibel (Guest/Customer):** Pencarian produk menggunakan kata kunci dinamis berbasis MongoDB Regex (Case-Insensitive).
- **Pembayaran Midtrans(Customer):** Pembayaran menggunakan Midtrans 
- **Sistem Pesanan Terintegrasi:** 
  - *Customer* dapat melakukan checkout dan menyelesaikan pesanan dengan pembayaran.
  - *Admin* dapat melakukan update user to seller 
  - *Seller* dapat memantau riwayat pesanan yang masuk (*ongoing* atau *completed*).

---

## 🛠️ Tech Stack

| Teknologi | Penggunaan |
| --- | --- |
| **Go (Golang)** | Bahasa pemrograman utama |
| **Gin Gonic** | HTTP Web Framework (Routing & Middleware) |
| **MongoDB Go Driver v2** | Database NoSQL untuk penyimpanan data |
| **JWT (Go-JWT)** | Pengamanan token autentikasi |
| **Midtrans** | Simulasi Pembayaran|

---

## ⚙️ Cara Menjalankan Projek di Lokal

Ikuti langkah-langkah berikut untuk menjalankan projek ini di komputer kamu:

### 1. Prasyarat
- Sudah menginstal [Go](https://go.dev/doc/install) (versi 1.20 ke atas)
- Sudah menginstal [MongoDB](https://www.mongodb.com/try/download/community) atau memiliki akun MongoDB Atlas

### 2. Kloning Repositori
```bash
git clone [https://github.com/UsernameKamu/nama-repo-kamu.git](https://github.com/UsernameKamu/nama-repo-kamu.git)
cd nama-repo-kamu