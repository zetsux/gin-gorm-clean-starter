# Gin GORM Template (Clean Architecture)

## Description

A Clean Architecture starter template for projects using Gin (Golang) and PostgreSQL, with GORM as the ORM.

## Architecture

```
/backend
│
├── /api
│   └── /v1
│       ├── /controller
│       │   ├── user.go
│       │   └── etc
│       └── /route
│           ├── user.go
│           └── etc
│
├── /common
│   ├── constants.go
│   ├── model.go
│   ├── response.go
│   └── etc
│
├── /config
│   ├── db.go
│   └── etc
│
├── /internal
│   ├── /dto
│   │   ├── user.go
│   │   └── etc
│   ├── /entity
│   │   └── user.go
│   │   └── etc
│   ├── /middleware
│   │   └── authentication.go
│   │   └── cors.go
│   │   └── authorization.go
│   ├── /repository
│   │   └── user.go
│   │   └── etc
│   └── /service
│       ├── user.go
│       └── etc
│
├── /migration
│   ├── /seeder
│   │   └── user.go
│   │   └── etc
│   └── migrator.go
│
├── /utils
│   ├── bcrypt.go
│   └── etc
│
└── main.go
```

### Explanation

- `/api/v1` : Directory yang berisi berbagai hal yang berkaitan dengan API seperti daftar endpoint yang disediakan (route) serta handler (controller) dari setiap endpointnya. Subdirectory `/v1` sendiri digunakan untuk menyimpan beberapa versi dari API yang dapat ditambahkan sesuai kebutuhan.

  - `/controller` : Directory untuk menyimpan hal-hal terkait dengan Controller yang merupakan bagian dari program yang berfungsi menerima Request dan memberikan Response.
  - `/route` : Directory untuk menyimpan hal-hal yang terkait dengan routing. Berisikan routes atau endpoints yang didukung beserta dengan metode Requestnya.

- `/common` : Directory yang berisi berbagai hal umum untuk digunakan di seluruh directory seperti Response, Struktur Model, dan Konstanta.

- `/config` : Directory yang berisi hal terkait konfigurasi aplikasi. Contohnya seperti konfigurasi Database.

- `/internal` : Directory yang berisi berbagai hal yang berkaitan dengan internal dari backend. Meliputi business logic, entitas, interaksi dengan database, dan lain-lain.

  - `/dto` : Directory untuk menyimpan DTO (Data Transfer Object) adalah placeholder untuk suatu object lain yang digunakan untuk menampung data Request dan Response.
  - `/entity` : Directory untuk menyimpan entitas atau model yang digunakan baik di migrasi maupun di aplikasi.
  - `/middleware` : Directory untuk menyimpan Middleware yang merupakan penengah dari suatu operasi.
  - `/repository` : Directory untuk menyimpan hal-hal terkait Repository yang merupakan lapisan yang berhubungan langsung dengan Database.
  - `/service` : Directory untuk menyimpan hal-hal terkait Service yang merupakan lapisan yang bertanggung jawab menangani alur atau logika bisnis aplikasi.

- `/migration`: Directory untuk menyimpan hal-hal terkait migrasi dan juga seeding terhadap Database.

  - `/seeder` : Directory untuk menyimpan hal-hal yang diperlukan untuk seeding terhadap Database dengan dipisahkan sesuai entitas.

- `/utils` : Directory untuk kode terkait fungsi-fungsi utilitas atau pembantu lainnya yang bisa digunakan di berbagai directory lainnya.

## Pre-requisites

1. Create the database in PostgreSQL with the name equal to the value of DB_NAME in `.env`
2. Use the command `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";` on the database terminal

## How to Run?

1. Use the command `make run` (or use `go run main.go` instead, if `make` is unable to be used) to run the application

## API Documentation (Postman)

Link : https://documenter.getpostman.com/view/25087235/2s9YXfcizj
