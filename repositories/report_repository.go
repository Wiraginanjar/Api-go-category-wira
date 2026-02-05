package repositories

import (
	"database/sql"
	"tugas-go/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetDailyReport() (*models.DailyReport, error) {
	report := &models.DailyReport{}

	// 1. Query Total Revenue & Total Transaksi Hari Ini
	// Menggunakan COALESCE agar jika null (tidak ada data) kembali 0
	queryTotals := `
		SELECT 
			COALESCE(SUM(total_amount), 0), 
			COUNT(id) 
		FROM transactions 
		WHERE DATE(created_at) = CURRENT_DATE
	`
	err := repo.db.QueryRow(queryTotals).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// 2. Query Produk Terlaris Hari Ini
	// Join transaction_details ke transactions untuk filter tanggal, dan ke product untuk ambil nama
	queryBestSeller := `
		SELECT 
			p.name, 
			SUM(td.quantity) as total_qty
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN product p ON td.product_id = p.id
		WHERE DATE(t.created_at) = CURRENT_DATE
		GROUP BY p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`
	
	var bestSeller models.BestSellingProduct
	err = repo.db.QueryRow(queryBestSeller).Scan(&bestSeller.Nama, &bestSeller.QtyTerjual)
	
	if err == sql.ErrNoRows {
		// Jika tidak ada penjualan hari ini, set null atau object kosong
		report.ProdukTerlaris = nil 
	} else if err != nil {
		return nil, err
	} else {
		report.ProdukTerlaris = &bestSeller
	}

	return report, nil
}

func (repo *ReportRepository) GetReportByDateRange(startDate, endDate string) (*models.DailyReport, error) {
	report := &models.DailyReport{}

	// 1. Query Total Revenue & Transaksi dalam rentang tanggal
	// Menggunakan parameter $1 dan $2 untuk tanggal
	queryTotals := `
		SELECT 
			COALESCE(SUM(total_amount), 0), 
			COUNT(id) 
		FROM transactions 
		WHERE DATE(created_at) >= $1 AND DATE(created_at) <= $2
	`
	
	err := repo.db.QueryRow(queryTotals, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// 2. Query Produk Terlaris dalam rentang tanggal
	queryBestSeller := `
		SELECT 
			p.name, 
			SUM(td.quantity) as total_qty
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN product p ON td.product_id = p.id
		WHERE DATE(t.created_at) >= $1 AND DATE(t.created_at) <= $2
		GROUP BY p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`

	var bestSeller models.BestSellingProduct
	err = repo.db.QueryRow(queryBestSeller, startDate, endDate).Scan(&bestSeller.Nama, &bestSeller.QtyTerjual)

	if err == sql.ErrNoRows {
		report.ProdukTerlaris = nil
	} else if err != nil {
		return nil, err
	} else {
		report.ProdukTerlaris = &bestSeller
	}

	return report, nil
}