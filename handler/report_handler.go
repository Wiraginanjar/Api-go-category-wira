package handlers

import (
	"net/http"
	"encoding/json"
	"tugas-go/repositories"
)

type ReportHandler struct {
	repo *repositories.ReportRepository
}

func NewReportHandler(repo *repositories.ReportRepository) *ReportHandler {
	return &ReportHandler{repo: repo}
}

func (h *ReportHandler) GetDailyReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.repo.GetDailyReport()
	if err != nil {
		http.Error(w, "Gagal mengambil data report: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *ReportHandler) GetReportByDate(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter dari URL
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// Validasi sederhana: Pastikan parameter tidak kosong
	if startDate == "" || endDate == "" {
		http.Error(w, "Parameter 'start_date' dan 'end_date' wajib diisi (Format: YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	// Panggil Repository
	report, err := h.repo.GetReportByDateRange(startDate, endDate)
	if err != nil {
		http.Error(w, "Gagal mengambil data report: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}