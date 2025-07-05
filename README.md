# 🧾 koer-tax-service

`koer-tax-service` adalah layanan microservice berbasis **Go**, menggunakan **gRPC**, **gRPC-Gateway (HTTP/REST)**, serta menerapkan **Clean Architecture** dan **separation of concerns** yang ketat. Setiap layer dapat diuji secara unit (unit testable), cocok untuk pengembangan skala besar dan enterprise.

---

## 🚀 Fitur Utama

- ✨ Clean Architecture dengan pemisahan tanggung jawab (separation of concerns)
- ⚙️ Dual API support: gRPC + HTTP/REST (via gRPC-Gateway)
- ✅ Setiap layer dapat di-*unit test*
- 🔒 Validasi input yang robust
- 📄 Logging terstruktur
- 🧪 Struktur testing yang modular

---

## 🧱 Struktur Proyek
koer-tax-service/
├── assets/ # File statis (contoh: font, img, dll)
├── cmd/
│ └── api/
│ └── main.go # Entry-point aplikasi
├── config/ # Konfigurasi aplikasi ( env loader, dsb)
├── pkg/ # Paket utilitas umum
│ ├── constants/ # Definisi konstanta global
│ ├── logs/ # Logging terstruktur
│ ├── utils/ # Fungsi bantu/utilitas umum
│ └── validator/ # Validasi input dan payload
├── proto/ # File .proto API service
├── protogen/ # Hasil generate dari buf/grpc (pb.go, gw.go, dsb)
├── server/
│ ├── domain/ # Entity & contract interface (model + interface)
│ ├── handler/ # Handler gRPC & REST (HTTP Gateway)
│ ├── infrastructure/ # Setup koneksi eksternal (DB, Redis, dll)
│ ├── repository/ # Implementasi akses data (PostgreSQL, dsb)
│ └── usecase/ # Bisnis logic / interactor
├── tests/ # Unit & integration test
└── go.mod / go.sum # Dependency & module management

---

## 🛠️ Build & Jalankan

### 1. Clone Repo

Kamu bisa menggunakan repo ini sebagai template atau mengkloning secara langsung:

#### 🔁 Opsi 1: Gunakan sebagai Template
Klik tombol **"Use this template"** di GitHub untuk membuat repo baru berdasarkan ini.

#### 📥 Opsi 2: Clone Manual

```bash
git clone https://github.com/kurnhyalcantara/koer-tax-service.git
cd koer-tax-service
```
---

## 🔁 Inisialisasi dan Update Submodule Proto

proto/ adalah Git submodule yang berisi definisi .proto dari seluruh layanan.

### 1. Pertama Kali Clone

```bash
git submodule update --init --recursive
```

### 2. Menarik Update Terbaru dari Repo Submodule

```bash
git submodule update --remote
```
---

### Generate gRPC & HTTP Gateway

Gunakan skrip generate.sh di dalam folder proto/ untuk meng-generate file .pb.go, .grpc.pb.go, dan .gw.go

```bash
./generate.sh <service-name>
```

### Running the Service

```bash
make run
```

### ✅ Testing
```bash
go test -v ./..
```

### 🧪 Arsitektur Clean

Handler (gRPC / REST)
       ↓
    Usecase (Business Logic)
       ↓
Repository Interface ← Repository Impl (PostgreSQL, dll)
       ↑
     Domain (Model + Interface Contract)

### 📜 Lisensi

MIT License © 2025 — Kurniawan
