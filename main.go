package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"tugas-go/database"
	"tugas-go/handler"
	"tugas-go/middleware"
	"tugas-go/repositories"
	"tugas-go/services"
	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
	API_KEY string `mapstructure:"API_KEY"`
}
func main() {
	
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
		API_KEY: viper.GetString("API_KEY"),
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	apiKeyMiddleware := middleware.APIKey(config.API_KEY)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	reportRepo := repositories.NewReportRepository(db)
	reportHandler := handlers.NewReportHandler(reportRepo)

	http.HandleFunc("/api/kategori", middleware.CORS(middleware.Logger(categoryHandler.HandleCategory)))
	http.HandleFunc("/api/kategori/", middleware.CORS(middleware.Logger(apiKeyMiddleware(categoryHandler.HandleCategoryByID))))
	http.HandleFunc("/api/produk", middleware.CORS(middleware.Logger(productHandler.HandleProducts)))
	http.HandleFunc("/api/produk/", middleware.CORS(middleware.Logger(apiKeyMiddleware(productHandler.HandleProductByID))))
	http.HandleFunc("/api/checkout", middleware.CORS(middleware.Logger(apiKeyMiddleware(transactionHandler.HandleCheckout)))) // POST
	http.HandleFunc("/api/report/hari-ini", middleware.CORS(middleware.Logger(reportHandler.GetDailyReport)))
	http.HandleFunc("/api/report", middleware.CORS(middleware.Logger(reportHandler.GetReportByDate)))
	http.HandleFunc("/api/health", middleware.CORS(middleware.Logger(categoryHandler.HealthCheck)))


	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("gagal running server", err)
	}
}
