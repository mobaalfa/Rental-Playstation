package konfigurasi

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func KoneksiDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/rental_playstation")
	if err != nil {
		log.Fatal("Gagal membuka koneksi:", err)
	}
	
	err = db.Ping()
	if err != nil {
		log.Fatal("Gagal terkoneksi ke database:", err)
	}

	return db
}
