# Gin GORM Template (Clean Architecture)

## Description

A starter template for projects using Gin or GORM (Golang)

## API Documentation (Postman)

Link : https://documenter.getpostman.com/view/25087235/2s9YXfcizj

## Pre-requisites

1. Create the database in PostgreSQL with the name equal to the value of DB_NAME in `.env`
2. Use the command `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";` on the database terminal

## Architecture Explanation

- common : Directory yang berisi berbagai hal untuk digunakan di seluruh directory seperti Response, Struktur Model, dan Konstanta.
- config : Directory yang berisi hal terkait konfigurasi seperti konfigurasi Database.
- controller : Directory untuk menyimpan hal-hal terkait dengan Controller yang merupakan bagian dari program yang berfungsi menerima Request dan memberikan Response.
- dto : Directory untuk menyimpan DTO (Data Transfer Object) adalah placeholder untuk suatu object lain yang digunakan untuk menampung data Request dan Response.
- entity : Directory untuk menyimpan entitas atau model yang digunakan.
- middleware : Directory untuk menyimpan Middleware yang merupakan penengah dari suatu operasi.
- repository : Directory untuk menyimpan hal-hal terkait Repository yang merupakan lapisan yang berhubungan langsung dengan Database.
- routes : Directory untuk menyimpan hal-hal yang terkait dengan routing. Berisikan routes atau endpoints yang didukung beserta dengan metode Requestnya.
- seeder : Directory untuk menyimpan hal-hal terkait migrasi dan seeding terhadap Database.
- utils : Directory untuk kode terkait fungsi-fungsi utilitas atau pembantu lainnya yang bisa digunakan di berbagai directory lainnya.
