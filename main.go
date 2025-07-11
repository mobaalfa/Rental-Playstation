package main

import (
	"net/http"

	"rental-playstation/kontrol"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/aset/").Handler(http.StripPrefix("/aset/", http.FileServer(http.Dir("aset"))))

	// Routing
	r.HandleFunc("/", kontrol.Dashboard).Methods("GET")
	r.HandleFunc("/penyewa", kontrol.TampilPenyewa).Methods("GET")
	r.HandleFunc("/tambah", kontrol.FormTambahPenyewa).Methods("GET")
	r.HandleFunc("/tambah", kontrol.ProsesTambahPenyewa).Methods("POST")
	r.HandleFunc("/hapus/{id}", kontrol.HapusPenyewa).Methods("GET")
	r.HandleFunc("/laporan", kontrol.CetakPDF).Methods("GET")

	r.HandleFunc("/edit/{id}", kontrol.TampilkanFormEdit).Methods("GET")
	r.HandleFunc("/edit/{id}", kontrol.ProsesEditPenyewa).Methods("POST")

	http.ListenAndServe(":8080", r)
}
