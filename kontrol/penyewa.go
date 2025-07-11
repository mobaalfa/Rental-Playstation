package kontrol

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	"rental-playstation/konfigurasi"
	"rental-playstation/model"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("tampilan/layout.html", "tampilan/dashboard.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Title": "Dashboard",
	})
}

func TampilPenyewa(w http.ResponseWriter, r *http.Request) {
	db := konfigurasi.KoneksiDB()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM penyewa")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data []model.Penyewa
	for rows.Next() {
		var p model.Penyewa
		rows.Scan(&p.ID, &p.NamaPenyewa, &p.JenisPS, &p.TanggalSewa, &p.TanggalKembali, &p.Harga)
		data = append(data, p)
	}

	tmpl, err := template.ParseFiles("tampilan/layout.html", "tampilan/daftar.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Title": "Daftar Penyewa",
		"Data":  data,
	})
}

func FormTambahPenyewa(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("tampilan/layout.html", "tampilan/form_tambah.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Title": "Tambah Penyewa",
	})
}

func ProsesTambahPenyewa(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nama := r.FormValue("nama_penyewa")
	jenis := r.FormValue("jenis_playstation")
	sewa := r.FormValue("tanggal_sewa")
	kembali := r.FormValue("tanggal_kembali")
	harga, _ := strconv.Atoi(r.FormValue("harga"))

	db := konfigurasi.KoneksiDB()
	defer db.Close()

	_, err := db.Exec("INSERT INTO penyewa (nama_penyewa, jenis_playstation, tanggal_sewa, tanggal_kembali, harga) VALUES (?, ?, ?, ?, ?)",
		nama, jenis, sewa, kembali, harga)
	if err != nil {
		http.Error(w, "Gagal menyimpan data", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/penyewa", http.StatusSeeOther)
}

func HapusPenyewa(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	db := konfigurasi.KoneksiDB()
	defer db.Close()

	db.Exec("DELETE FROM penyewa WHERE id = ?", id)
	http.Redirect(w, r, "/penyewa", http.StatusSeeOther)
}

func CetakPDF(w http.ResponseWriter, r *http.Request) {
	db := konfigurasi.KoneksiDB()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM penyewa")
	if err != nil {
		http.Error(w, "Gagal ambil data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var penyewas []model.Penyewa
	for rows.Next() {
		var p model.Penyewa
		rows.Scan(&p.ID, &p.NamaPenyewa, &p.JenisPS, &p.TanggalSewa, &p.TanggalKembali, &p.Harga)
		penyewas = append(penyewas, p)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(214, 40, 40) // Merah terang
	pdf.CellFormat(190, 10, "LAPORAN PENYEWA PLAYSTATION", "0", 1, "C", false, 0, "")
	pdf.Ln(6)

	// Table Header
	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(64, 1, 17)     // Merah gelap
	pdf.SetTextColor(255, 255, 255) // Putih
	headers := []string{"ID", "Nama", "PS", "Tgl Sewa", "Tgl Kembali", "Harga"}
	widths := []float64{10, 45, 20, 30, 30, 30}

	for i, h := range headers {
		pdf.CellFormat(widths[i], 10, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table Body
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(245, 245, 245)
	fill := false

	for _, p := range penyewas {
		pdf.CellFormat(widths[0], 9, fmt.Sprintf("%d", p.ID), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(widths[1], 9, p.NamaPenyewa, "1", 0, "", fill, 0, "")
		pdf.CellFormat(widths[2], 9, p.JenisPS, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(widths[3], 9, p.TanggalSewa, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(widths[4], 9, p.TanggalKembali, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(widths[5], 9, fmt.Sprintf("Rp%d", p.Harga), "1", 0, "R", fill, 0, "")
		pdf.Ln(-1)
		fill = !fill
	}

	w.Header().Set("Content-Type", "application/pdf")
	err = pdf.Output(w)
	if err != nil {
		http.Error(w, "Gagal generate PDF", http.StatusInternalServerError)
	}
}

func TampilkanFormEdit(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	db := konfigurasi.KoneksiDB()
	defer db.Close()

	var p model.Penyewa
	err := db.QueryRow("SELECT * FROM penyewa WHERE id = ?", id).Scan(
		&p.ID, &p.NamaPenyewa, &p.JenisPS, &p.TanggalSewa, &p.TanggalKembali, &p.Harga)
	if err != nil {
		http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("tampilan/layout.html", "tampilan/form_edit.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Title": "Edit Penyewa",
		"Data":  p,
	})
}

func ProsesEditPenyewa(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := mux.Vars(r)["id"]
	nama := r.FormValue("nama_penyewa")
	jenis := r.FormValue("jenis_playstation")
	sewa := r.FormValue("tanggal_sewa")
	kembali := r.FormValue("tanggal_kembali")
	harga, _ := strconv.Atoi(r.FormValue("harga"))

	db := konfigurasi.KoneksiDB()
	defer db.Close()

	_, err := db.Exec(`UPDATE penyewa SET 
		nama_penyewa=?, jenis_playstation=?, tanggal_sewa=?, tanggal_kembali=?, harga=? 
		WHERE id=?`,
		nama, jenis, sewa, kembali, harga, id)

	if err != nil {
		http.Error(w, "Gagal update data", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/penyewa", http.StatusSeeOther)
}
