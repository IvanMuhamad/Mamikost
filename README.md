Ini merupakan API berbasis Go untuk mengelola pemesanan kamar kost. Sistem ini memungkinkan pengguna untuk membuat akun, login kedalam akun sekaligus logout, menambah dan mengkofigurasikan kategori kost, menambah dan mengkofigurasikan properti kamar kost, mengupload foto kamar kost, menambahkan item kedalam cart, melakukan pemesanan, dan melihat detail pesanan.
API ini menggunakan Gin sebagai framework web dan SQLC untuk database. Keamanan berbasis token diterapkan untuk semua route yang dilindungi.

Tools yang digunakan:

Go: Bahasa pemrograman utama.

Gin: Framework web HTTP.

SQLC: Alat untuk menghasilkan kode Go dari query SQL.

PostgreSQL: Database untuk menyimpan informasi properti kost dan pesanan.

JWT: Autentikasi berbasis token.

Swagger: Dokumentasi API.

Postman: Pengujian API

Docker: Containerisasi untuk memudahkan deployment dan pengujian.

