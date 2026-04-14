package main

import (
	"log"
	"os"
	
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	
	"faturamento-service/internal/client"
	"faturamento-service/internal/handler"
	"faturamento-service/internal/models"
	"faturamento-service/internal/repository"
	"faturamento-service/internal/routes"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env não encontrado, usando variáveis do sistema")
	}
	
	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=postgres user=korp password=korp123 dbname=faturamento_db port=5432 sslmode=disable"
	}
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Erro ao conectar banco:", err)
	}
	
	// Auto migrate
	db.AutoMigrate(&models.NotaFiscal{}, &models.NotaItem{})
	
	// Initialize dependencies
	notaRepo := repository.NewNotaRepository(db)
	estoqueClient := client.NewEstoqueClient()
	notaHandler := handler.NewNotaHandler(notaRepo, estoqueClient)
	
	// Setup router
	r := gin.Default()
	routes.SetupRoutes(r, notaHandler)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	
	log.Printf("🚀 faturamento-service iniciado na porta %s", port)
	r.Run(":" + port)
}
