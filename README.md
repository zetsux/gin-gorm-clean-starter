# Gin GORM Template (Clean Architecture)

## Contents

- [Description](#description)
- [Architecture](#architecture)

  - [Explanation (EN)](#explanation-en)
  - [Explanation (ID)](#explanation-id)

- [Pre-requisites](#pre-requisites)

  - [PostgreSQL Requirements](#postgresql-requirements)
  - [GitHooks Requirements](#githooks-requirements)

- [How to Run?](#how-to-run)
- [API Documentation (Postman)](#api-documentation-postman)

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
│   ├── /middleware
│   │   └── authentication.go
│   │   └── cors.go
│   │   └── authorization.go
│   │   └── etc
│   ├── /standard
│   │   └── constants.go
│   │   └── errors.go
│   │   └── model.go
│   │   └── request.go
│   │   └── response.go
│   │   └── etc
│   └── /utils
│       └── bcrypt.go
│       └── file.go
│       └── etc
│
├── /config
│   ├── db.go
│   └── etc
│
├── /database
│   ├── /seeder
│   │   └── user.go
│   │   └── etc
│   └── migrator.go
│
├── /internal
│   ├── /dto
│   │   ├── user.go
│   │   └── etc
│   ├── /entity
│   │   └── user.go
│   │   └── etc
│   ├── /repository
│   │   └── user.go
│   │   └── etc
│   └── /service
│       ├── user.go
│       └── etc
│
└── main.go
```

### Explanation (EN)

- `/api/v1` : The directory for things related to API like all available endpoints (route) and the handlers for each endpoints (controller). Subdirectory `/v1` is used for easy version control in case of several development phase.

  - `/controller` : The directory for things related to the Controller layer which is the part of program that handle requests and return responses.
  - `/route` : The directory for things related with routing. Therefore filled with every available supported routes / endpoints along with the request method and used middleware.

- `/common` : The directory for common things that are frequently used all over the architectures.

  - `/middleware` : The directory for Middlewares which are mechanism that intercept a HTTP request and response process before handled directly by the controller of an endpoint.
  - `/standard` : The directory for standard things such as variables, constants, and functions to be used in other directories. It consists of things like response, request, error, model base structure, constants, etc.
  - `/utils` : The directory to store utility / helper functions that can be used in other directories.

- `/config` : The directory for things related to program configuration like database configuration.

- `/database`: The directory for things related to the database for example migrations and seeders.

  - `/seeder` : The directory for things related to database seeding for each entity.

- `/internal` : The directory for things related to the internal side of the Back End. It consists of things like business logic, entities, and database interaction.

  - `/dto` : The directory to store DTO (Data Transfer Object) which is a placeholder for other objects, mainly to store data for requests and responses. Also stores error and messages for each entity.
  - `/entity` : The directory for things related to entities / models which are available on the database via migration that are represented by structs.
  - `/repository` : The directory for things related to the Repository layer which is the layer that is responsible to interact directly with the database.
  - `/service` : The directory for things related to the Service layer which is the layer that is responsible for the flow / business logic of the app.

### Explanation (ID)

- `/api/v1` : Directory yang berisi berbagai hal yang berkaitan dengan API seperti daftar endpoint yang disediakan (route) serta handler (controller) dari setiap endpointnya. Subdirectory `/v1` sendiri digunakan untuk menyimpan beberapa versi dari API yang dapat ditambahkan sesuai kebutuhan.

  - `/controller` : Directory untuk menyimpan hal-hal terkait dengan Controller yang merupakan bagian dari program yang berfungsi menerima request dan memberikan response.
  - `/route` : Directory untuk menyimpan hal-hal yang terkait dengan routing. Berisikan routes atau endpoints yang didukung beserta dengan metode requestnya.

- `/common` : Directory yang berisi berbagai hal umum untuk digunakan di seluruh directory.

  - `/middleware` : Directory untuk menyimpan Middleware yang merupakan mekanisme yang menengahi proses HTTP request dan response sebelum ditangani secara langsung oleh controller setiap route.
  - `/standard` : Directory yang berisi berbagai variabel, konstanta, maupun fungsi standar untuk digunakan di berbagai directory lainnya seperti response, request, error, struktur dasar model, konstanta, dan lain-lain.
  - `/utils` : Directory untuk kode terkait fungsi-fungsi utilitas atau pembantu lainnya yang bisa digunakan di berbagai directory lainnya.

- `/config` : Directory yang berisi hal terkait konfigurasi aplikasi. Contohnya seperti konfigurasi database.

- `/database`: Directory untuk menyimpan hal-hal terkait migrasi dan juga seeding terhadap database.

  - `/seeder` : Directory untuk menyimpan hal-hal yang diperlukan untuk seeding terhadap database dengan dipisahkan sesuai entitas.

- `/internal` : Directory yang berisi berbagai hal yang berkaitan dengan sisi internal dari Back End. Meliputi business logic, entitas, maupun interaksi dengan database.

  - `/dto` : Directory untuk menyimpan DTO (Data Transfer Object) adalah placeholder untuk suatu object lain yang digunakan untuk menampung data request dan response. Juga menyimpan error serta pesan yang spesifik untuk setiap entitas.
  - `/entity` : Directory untuk menyimpan entitas atau model yang digunakan baik di migrasi maupun di aplikasi.
  - `/repository` : Directory untuk menyimpan hal-hal terkait Repository yang merupakan lapisan yang berhubungan langsung dengan database.
  - `/service` : Directory untuk menyimpan hal-hal terkait Service yang merupakan lapisan yang bertanggung jawab menangani alur atau logika bisnis aplikasi.

## Pre-requisites

### PostgreSQL Requirements

1. Create the database in PostgreSQL with the name equal to the value of DB_NAME in `.env`
2. Use the command `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";` on the database terminal

### GitHooks Requirements

> Note : GitHooks is not mandatory for this starter. Only do the steps below if you want to apply & use it.

1. Install golangci-lint as the linters aggregator for pre-commit linting by executing `go install github.com/golangci/golangci-lint@latest`. Alternatively, you can follow the recommended method, which involves installing the binary from the [official source](https://golangci-lint.run/usage/install/#binaries)
2. Install commitlint as the conventional commit message checker by executing `go install github.com/conventionalcommit/commitlint@latest`. Alternatively, you can follow the recommended method, which involves installing the binary from the [official source](https://github.com/conventionalcommit/commitlint/releases)
3. Configure your git's hooks path to be linked to the `.githooks` directory on this repository by executing `git config core.hooksPath .githooks`

## How to Run?

1. Use the command `make tidy` (or use `go mod tidy` instead, if `make` is unable to be used) to adjust the dependencies accordingly
2. Use the command `make run` (or use `go run main.go` instead, if `make` is unable to be used) to run the application

## API Documentation (Postman)

Link : https://documenter.getpostman.com/view/25087235/2s9YXfcizj
